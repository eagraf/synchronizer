package messenger

import (
	"encoding/json"
	"testing"
	"time"
)

// MessageBuilder tests

// Unit test MessageBuilder.NewMessage
func TestNewMessage(t *testing.T) {
	// Call function
	mb := MessageBuilder{}
	mb.NewMessage("test-message", "test-request-id")
	m := mb.message

	// Check valid uuid
	if len(m.metadata.UUID) != 36 {
		t.Error("Invalid UUID")
	}

	// Check timestamp is
	if _, err := time.Parse(time.StampMilli, m.metadata.Timestamp); err != nil {
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
	if m.metadata.Headers == nil {
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
	if m.offset != 0 { // TODO make sure this is correct
		t.Error("Incorrect offset")
	}
}

func TestAddHeaderString(t *testing.T) {
	mb := MessageBuilder{}
	mb.NewMessage("test-message", "test-request-id").
		AddHeader("test-key", "test-value")

	v, ok := mb.message.metadata.Headers["test-key"]
	if ok != true {
		t.Error("Key is not added")
	}
	if v != "test-request-id" {
		t.Error("Value incorrect")
	}
}

func TestAddHeaderJSON(t *testing.T) {
	mb := MessageBuilder{}
	headerMap := make(map[string]interface{})
	headerMap["a"] = "b"
	mb.NewMessage("test-message", "test-request-id").
		AddHeader("test-key", headerMap)

	var val map[string]interface{}
	v, ok := mb.message.metadata.Headers["test-key"]
	val = v.(map[string]interface{})

	if ok != true {
		t.Error("Key is not added")
	}
	if len(val) != 1 {
		t.Error("Length of map incorrect")
	}
	if val["a"] != "b" {
		t.Error("Value incorrect")
	}
}

func TestAddHeaderStruct(t *testing.T) {
	type testStruct struct {
		a int
	}
	headerStruct := testStruct{1}

	mb := MessageBuilder{}
	_, err := mb.NewMessage("test-message", "test-request-id").
		AddHeader("test-key", headerStruct).
		Done()

	if err.Error() != "Invalid header value" {
		t.Error("Incorrect error value")
	}

	// Test that the header was not added
	if len(mb.message.metadata.Headers) != 0 {
		t.Error("Header was added")
	}
}

func TestSetPayload(t *testing.T) {
	payload := []byte("Hello, World")
	mb := MessageBuilder{}
	mb.NewMessage("test-message", "test-request-id").
		SetPayload(payload)

	if mb.message.payload == nil {
		t.Error("Payload was not added")
	}
	if len(mb.message.payload) != 12 {
		t.Error("Payload incorrect length")
	}
}

func TestSetPayloadEmpty(t *testing.T) {
	payload := make([]byte, 0)
	mb := MessageBuilder{}
	_, err := mb.NewMessage("test-message", "test-request-id").
		SetPayload(payload).
		Done()
	if err.Error() != "Payload is empty" {
		t.Error("Incorrect error value")
	}
}

func TestDone(t *testing.T) {
	payload := []byte("Hello, World")
	mb := MessageBuilder{}
	message, err := mb.NewMessage("test-message", "test-request-id").
		AddHeader("test-key", "test-value").
		SetPayload(payload).
		Done()

	if err != nil {
		t.Error("No error value expected")
	}

	if message.done == false {
		t.Error("message.done not set")
	}

	// message should delete reference to payload after compression for garbage collection
	if len(message.payload) > 0 {
		t.Error("Payload length not 0")
	}

	// TODO stronger size check
	if len(message.deflated) == 0 {
		t.Error("Message was not compressed")
	}

	marshalledMetadata, _ := json.Marshal(message.metadata)
	if message.offset != int32(len(marshalledMetadata)+4) {
		t.Error("Offset value incorrect")
	}

	if message.received == true {
		t.Error("message.received is incorrect")
	}

	if mb.message != nil {
		t.Error("mb.message != nil")
	}

}
func TestModifyAfterDone(t *testing.T) {
	mb := MessageBuilder{}
	mb.NewMessage("test-message", "test-request-id").Done()
	nmb := mb.AddHeader("test-key", "test-value")
	// Defined behaviour: if AddHeader or SetPayload are called after Done, they return nil
	if nmb != nil {
		t.Error("AddHeader does not return nil")
	}

	nmb = mb.SetPayload(make([]byte, 3))
	if nmb != nil {
		t.Error("SetPayload does not return nil")
	}
}

func TestFromBuffer(t *testing.T) {
	payload := []byte("Hello, World")
	headerMap := make(map[string]interface{})
	headerMap["a"] = "b"
	mb := MessageBuilder{}
	message, _ := mb.NewMessage("test-message", "test-request-id").
		AddHeader("test-key", "test-value").
		AddHeader("test-key-2", headerMap).
		SetPayload(payload).
		Done()

	message2, err := mb.FromBuffer(message.deflated)
	if err != nil {
		t.Error("Error is not nil")
	}

	if message2.metadata.MessageType != "test-request-id" {
		t.Error("Invalid metadata")
	}
	if message2.metadata.Request != "test-request-id" {
		t.Error("Invalid metadata")
	}
	if message2.metadata.HasPayload != true {
		t.Error("HasPayload returns false")
	}
	if message2.metadata.Headers["test-key"] != "test-value" {
		t.Error("Header missing")
	}
	if (message2.metadata.Headers["test-key-2"]).(map[string]interface{})["a"] != "b" {
		t.Error("JSON header missing")
	}
}

func TestFromBufferError(t *testing.T) {
	mb := MessageBuilder{}
	_, err := mb.FromBuffer([]byte("Hello"))
	if err != nil {
		t.Error("FromBuffer should return an error")
	}
}

// Message Tests

func TestGetMetadata(t *testing.T) {
	mb := MessageBuilder{}
	message, _ := mb.NewMessage("test-message", "test-request-id").Done()

	if message.GetMetadata().MessageType != "test-request-id" {
		t.Error("Invalid metadata")
	}
	if message.GetMetadata().Request != "test-request-id" {
		t.Error("Invalid metadata")
	}
}

func TestGetPayload(t *testing.T) {
	p := []byte("Hello, World")
	mb := MessageBuilder{}
	message, _ := mb.NewMessage("test-message", "test-request-id").
		SetPayload(p).
		Done()

	payload, _ := message.GetPayload()
	if payload == nil {
		t.Error("Payload is nil")
	}
	if len(payload) != 12 {
		t.Error("Payload length incorrect")
	}
}

func TestGetEmptyPayload(t *testing.T) {
	mb := MessageBuilder{}
	message, _ := mb.NewMessage("test-message", "test-request-id").Done()

	_, err := message.GetPayload()
	if err != nil {
		t.Error("Error should be returned")
	}
}

func TestGetHeaderString(t *testing.T) {
	mb := MessageBuilder{}
	message, _ := mb.NewMessage("test-message", "test-request-id").
		AddHeader("test-key", "test-value").
		Done()

	if message.GetHeader("test-key") != "test-value" {
		t.Error("Incorrect value returned")
	}
}

func TestGetHeaderJSON(t *testing.T) {
	mb := MessageBuilder{}
	headerMap := make(map[string]interface{})
	headerMap["a"] = "b"
	message, _ := mb.NewMessage("test-message", "test-request-id").
		AddHeader("test-key", headerMap).
		Done()

	if message.GetHeader("test-key").(map[string]interface{})["a"] != "b" {
		t.Error("JSON header failed to get key")
	}
}
