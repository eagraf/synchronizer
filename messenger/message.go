package messenger

/*
 */

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
	offset   int32     // The byte length of serialized JSON + 4
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
}

// Message methods

// GetMetadata returns all message metadata
func (m *Message) GetMetadata() *Metadata {
	return nil
}

// GetHeader returns the value of a single header key
func (m *Message) GetHeader(key string) interface{} {
	return nil
}

// GetPayload returns the slice at the correct position for the deflated byte array
func (m *Message) GetPayload() ([]byte, error) {
	return nil, nil
}

// MessageBuilder methods

// FromBuffer creates message struct from received buffer, and decompress
func (mb *MessageBuilder) FromBuffer([]byte) *Message {
	return nil
}

// NewMessage begins message creation, taking metadata as input
func (mb *MessageBuilder) NewMessage(messageType string, request string) *Message {
	return nil
}

// AddHeader adds a key value pair to metadata. The value can be either a string, number, boolean, or JSON object
func (mb *MessageBuilder) AddHeader(key string, value interface{}) *Message {
	return nil
}

// SetPayload sets I/O data for workunit
func (mb *MessageBuilder) SetPayload([]byte) *Message {
	return nil
}

// Done completes message building and prepares for sending by compressing the data
func (mb *MessageBuilder) Done() *Message {
	return nil
}
