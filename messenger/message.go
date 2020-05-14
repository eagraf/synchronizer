package messenger

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
)

// Message encapsulates information passed over websocket connections
/*
 * Two ways of creating a Message
 *   (1) Building the message before sending it
 *   (2) Populating fields after the message has been received
 *
 * Message Structure:
 *  (1) Offset: [int32] Byte length of serialized JSON + 4.
 *  (2) Metadata: Serialized JSON with important information, and additional optional headers
 *  (3) Payload: Raw data that is I/O for a workunit. For large data, it is important that this is not copied repeatedly into new buffers
 *
 * Key idea: as large ammounts of data can be transferred, minimizing the number of memory allocs is important for performance
 */
type Message struct {
	offset   uint32    // The byte length of serialized JSON + 4
	metadata *Metadata // Important information about the message
	payload  []byte    // Compressed data
	deflated []byte    // TODO make sure this is freed after decompression (test it too?)
	received bool      // Was this message received or intended to be sent. received: T
	done     bool      // Utility for message builder to know when complete
}

// Metadata contains required information plus optional headers. It is sent as serialized JSON in messages
type Metadata struct {
	UUID        string                 `json:"uuid"`
	Timestamp   string                 `json:"timestamp"` // Use time formatted with time.StampMilli
	MessageType string                 `json:"messageType"`
	Request     string                 `json:"request"`
	HasPayload  bool                   `json:"hasPayload"`
	Headers     map[string]interface{} `json:"headers"`
}

// MessageBuilder constructs a new message struct
type MessageBuilder struct {
	message *Message
	err     error
}

// Message methods

// GetMetadata returns all message metadata
func (m *Message) GetMetadata() *Metadata {
	return m.metadata
}

// GetHeader returns the value of a single header key
func (m *Message) GetHeader(key string) (interface{}, bool) {
	header, ok := m.metadata.Headers[key]
	if !ok {
		return nil, false
	}
	return header, true
}

// GetPayload returns the slice at the correct position for the deflated byte array
func (m *Message) GetPayload() ([]byte, error) {
	if m.metadata.HasPayload == false || len(m.payload) == 0 {
		return nil, errors.New("This message has no payload")
	}
	return m.payload, nil
}

// FromBuffer creates message struct from received buffer, and decompress
func FromBuffer(buffer []byte) (*Message, error) {

	// Decompress message
	zr, err := zlib.NewReader(bytes.NewReader(buffer))
	if err != nil {
		return nil, err
	}
	// Read into byte array
	inflatedBuffer := new(bytes.Buffer)
	_, err = inflatedBuffer.ReadFrom(zr)
	if err != nil {
		return nil, err
	}
	inflated := inflatedBuffer.Bytes()

	// Determine offset
	offset := binary.LittleEndian.Uint32(inflated[0:4])

	// Get metadata
	var metadata *Metadata
	err = json.Unmarshal(inflated[4:offset], &metadata)
	fmt.Println(string(inflated[4:offset]))
	if err != nil {
		fmt.Println("hey!" + err.Error())
		return nil, err
	}

	// Construct message
	message := Message{
		offset:   offset,
		metadata: metadata,
		received: true,
		done:     true,
	}

	// Set payload if exists
	if metadata.HasPayload {
		message.payload = inflated[offset:]
	}

	return &message, nil
}

// MessageBuilder methods

// NewMessage begins message creation, taking metadata as input
func (mb *MessageBuilder) NewMessage(messageType string, request string) *MessageBuilder {

	metadata := Metadata{
		UUID:        uuid.New().String(),
		Timestamp:   time.Now().Format(time.StampMilli),
		MessageType: messageType,
		Request:     request,
		Headers:     make(map[string]interface{}),
		HasPayload:  false,
	}

	msg := Message{
		metadata: &metadata,
		received: false,
		done:     false,
	}
	mb.message = &msg

	return mb
}

// AddHeader adds a key value pair to metadata. The value can be either a string, number, boolean, or JSON object
func (mb *MessageBuilder) AddHeader(key string, value interface{}) *MessageBuilder {

	if mb.message == nil || mb.message.done {
		return nil
	}
	// Ensure type is correct
	// TODO potentially add support for more complex header types in the future
	// If key is already in headers, then replace
	if checkType(value) {
		mb.message.metadata.Headers[key] = value
	} else {
		// Add error to MessageBuilder
		mb.err = errors.New("Type " + reflect.TypeOf(value).String() + " not allowed in header.")
		return mb
	}

	return mb
}

// SetPayload sets I/O data for workunit
func (mb *MessageBuilder) SetPayload(payload []byte) *MessageBuilder {
	if mb.message == nil || mb.message.done {
		return nil
	}
	if len(payload) == 0 {
		mb.err = errors.New("Payload is empty")
		return mb
	}
	mb.message.payload = payload
	mb.message.metadata.HasPayload = true
	return mb
}

// Done completes message building and prepares for sending by compressing the data
// Returns an error if any of the previous build steps has been incorrect
func (mb *MessageBuilder) Done() (*Message, error) {
	// Check for error first
	if mb.err != nil {
		return nil, mb.err
	}

	// Marshal metadata
	metadataBuffer, err := json.Marshal(mb.message.metadata)
	if err != nil {
		return nil, err
	}
	mb.message.offset = uint32(len(metadataBuffer) + 4)

	mb.message.done = true
	m := mb.message
	mb.message = nil
	return m, nil
	// It makes a lot more sense to do this all at send time, so that a writer can connect all pieces without need for intermediate buffers
	/*

		// Message buffer
		buffer := make([]byte, 4 + len(metadataBuffer) + len(mb.message.payload))

		// Put offset
		binary.LittleEndian.PutUint32(buffer, mb.message.offset)

		// Put marshalled metadata
		copy(buffer[4:], metadataBuffer)

		// TODO investigate if there is a way to eliminate this copy
	*/
}

// Utility to determine valid header
func checkType(val interface{}) bool {
	switch val.(type) {
	case int:
		return true
	case int8:
		return true
	case int16:
		return true
	case int32:
		return true
	case int64:
		return true
	case uint:
		return true
	case uint8:
		return true
	case uint16:
		return true
	case uint32:
		return true
	case uint64:
		return true
	case float32:
		return true
	case float64:
		return true
	case string:
		return true
	case bool:
		return true
	default:
		return false
	}
}
