package messenger

import (
	"testing"
	"time"
)

// MessageBuilder tests

// Unit test MessageBuilder.NewMessage
func TestNewMessage(t *testing.T) {
	// Call function
	mb := MessageBuilder{}
	m := mb.NewMessage("test-message", "test-request-id")

	// Check valid uuid
	if len(m.metadata.UUID) != 36 {
		t.Error("Invalid UUID")
	}

	// Check timestamp is
	if _, err := time.Parse(time.StampMilli, m.metadata.Timestamp); err == nil {
		t.Error("Invalid timestamp")
	}

	// Check message type
	if m.metadata.MessageType != "test-message" {
		t.Error("Invalid message type")
	}

	// Check request
	if m.metadata.Request != "test-request-id" {
		t.Error("Invalid request id")
	}

	// Check headers are initialized
	if m.GetMetadata().Headers == nil {
		t.Error("Headers not initialized")
	}
	if len(m.metadata.Headers) != 0 {
		t.Error("There are more than zero headers")
	}

	// Check other values
	if m.payload != nil {
		t.Error("payload is not nil")
	}
	if m.deflated != nil {
		t.Error("deflated is not nil")
	}
	if m.received != false {
		t.Error("This message is being built, not received")
	}
	if m.offset != 4 { // TODO make sure this is correct
		t.Error("Incorrect offset")
	}
}

func TestAddHeaderString(t *testing.T) {

}

func TestAddHeaderJSON(t *testing.T) {

}

func TestAddHeaderStruct(t *testing.T) {

}

func TestSetPayload(t *testing.T) {

}

func TestSetPayloadEmpty(t *testing.T) {

}

func TestDone(t *testing.T) {

}

func TestDoneIncomplete(t *testing.T) {

}

// Message Tests

func TestGetMetadata(t *testing.T) {

}

func TestGetPayload(t *testing.T) {

}

func TestGetEmptyPayload(t *testing.T) {

}

func TestGetHeaderString(t *testing.T) {

}

func TestGetHeaderJSON(t *testing.T) {

}
