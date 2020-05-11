package messenger

/*
 * Message Structure:
 *  (1) Offset: [int32] Byte length of serialized JSON + 4.
 *  (2) Metadata: Serialized JSON with important information, and additional optional headers
 *  (3) Payload: Raw data that is I/O for a workunit. For large data, it is important that this is not copied repeatedly into new buffers
 */

// Internal representation of messages
// Key idea: as large ammounts of data can be transferred, minimizing the number of memory allocs is important for performance
type message struct {
	offset   int32     // The byte length of serialized JSON + 4
	metadata *Metadata // Important information about the message
	raw      []byte    // Compressed data
	deflated []byte    // TODO make sure this is freed after decompression (test it too?)
	received bool      // Was this message received or intended to be sent. received: T
	done     bool      // Utility for message builder to know when complete
}

// Metadata contains required information plus optional headers. It is sent as serialized JSON in messages
type Metadata struct {
	UUID        string                 `json:"uuid"`
	Timestamp   string                 `json:"timestamp"`
	MessageType string                 `json:"messageType"`
	Request     string                 `json:"request"`
	HasPayload  bool                   `json:"hasPayload"`
	Headers     map[string]interface{} `json:"headers"`
}

// Message encapsulates information passed over websocket connections
/*
 * Two ways of creating a Message
 *   (1) Building the message before sending it
 *   (2) Populating fields after the message has been received
 */
type _Message interface {
	GetMetaData() *Metadata
	GetPayload() []byte               // Importantly, this returns a slice indexed at the correct position within the _deflated_ byte array
	GetHeader(key string) interface{} // Return the value of a single header
}

// MessageBuilder constructs a new message struct
type MessageBuilder interface {
	FromBuffer([]byte) *Message                             // Create message struct from received buffer, and decompress
	NewMessage(messageType string, request string) *Message // Required metadata is included
	AddHeader(key string, value interface{}) *Message       // For adding optional metadata
	SetPayload([]byte) *Message                             // I/O for workunit
	Done() (*Message, error)                                // Prepares the message for sending by compressing it
}

// Internal representation of MessageBuilder, just contains incomplete message struct
type messageBuilder struct {
	message *message
}

// Message methods

// MessageBuilder methods

//func (mb *messageBuilder) FromBuffer([]byte) *Message {
// Unimplemented
