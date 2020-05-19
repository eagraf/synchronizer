package messenger

import (
	"errors"
	"time"
)

type messengerLog struct {
	openRequests      map[string]*RoundTrip
	completedRequests map[string]*RoundTrip // TODO this will grow out of control -> a better logging system is needed
}

// RoundTrip represents a request response pair
// TODO This might not encapsulate all possible messaging patterns, a more general graph based structure might be needed
type RoundTrip struct {
	requestID string
	request   *Message
	response  *Message
}

// Create a new instance of a messenger log
func newMessengerLog() *messengerLog {
	ml := messengerLog{
		openRequests:      make(map[string]*RoundTrip),
		completedRequests: make(map[string]*RoundTrip),
	}
	return &ml
}

// Duration gets the total duration of a request
func (rt *RoundTrip) Duration() (*time.Duration, error) {
	if rt.request == nil || rt.response == nil {
		return nil, errors.New("RoundTrip struct incomplete")
	}
	t1, err1 := time.Parse(time.StampMilli, rt.response.metadata.Timestamp)
	t2, err2 := time.Parse(time.StampMilli, rt.request.metadata.Timestamp)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	d := t2.Sub(t1)
	return &d, nil
}

func (l *messengerLog) GetIdentifier() string {
	// TODO universal identifier
	return "messenger_logger"
}

// OnReceive callback insert or complete roundtrip
func (l *messengerLog) OnReceive(topic string, message *Message) {
	// Check if it matches any active message request
	rt, ok := l.openRequests[message.metadata.Request]
	if ok == true {
		// Complete round trip
		rt.response = message
		l.completedRequests[rt.requestID] = rt
		delete(l.openRequests, rt.requestID)
	} else {
		// Otherwise create new request
		rt = new(RoundTrip)
		rt.requestID = message.metadata.Request
		rt.request = message
		l.openRequests[rt.requestID] = rt
	}
}

// OnSend callback insert or complete roundtrip
func (l *messengerLog) OnSend(topic string, message *Message) {
	// Check if it matches any active message request
	rt, ok := l.openRequests[message.metadata.Request]
	if ok == true {
		// Complete round trip
		rt.response = message
		l.completedRequests[rt.requestID] = rt
		delete(l.openRequests, rt.requestID)
	} else {
		// Otherwise create new request
		rt = new(RoundTrip)
		rt.requestID = message.metadata.Request
		rt.request = message
		l.openRequests[rt.requestID] = rt
	}
}

func (l *messengerLog) OnClose(topic string) {
	// pass
}
