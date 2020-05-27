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
	RequestID string
	Request   *Message
	Response  *Message
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
	if rt.Request == nil || rt.Response == nil {
		return nil, errors.New("RoundTrip struct incomplete")
	}
	t1, err1 := time.Parse(time.StampMilli, rt.Response.metadata.Timestamp)
	t2, err2 := time.Parse(time.StampMilli, rt.Request.metadata.Timestamp)
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}
	d := t2.Sub(t1)
	return &d, nil
}

func (ml *messengerLog) GetIdentifier() string {
	// TODO universal identifier
	return "messenger_logger"
}

// OnReceive callback insert or complete roundtrip
func (ml *messengerLog) OnReceive(topic string, message *Message) {
	// Check if it matches any active message request
	rt, ok := ml.openRequests[message.metadata.Request]
	if ok == true {
		// Complete round trip
		rt.Response = message
		ml.completedRequests[rt.RequestID] = rt
		delete(ml.openRequests, rt.RequestID)
	} else {
		// Otherwise create new request
		rt = new(RoundTrip)
		rt.RequestID = message.metadata.Request
		rt.Request = message
		ml.openRequests[rt.RequestID] = rt
	}
}

// OnSend callback insert or complete roundtrip
func (ml *messengerLog) OnSend(topic string, message *Message) {
	// Check if it matches any active message request
	rt, ok := ml.openRequests[message.metadata.Request]
	if ok == true {
		// Complete round trip
		rt.Response = message
		ml.completedRequests[rt.RequestID] = rt
		delete(ml.openRequests, rt.RequestID)
	} else {
		// Otherwise create new request
		rt = new(RoundTrip)
		rt.RequestID = message.metadata.Request
		rt.Request = message
		ml.openRequests[rt.RequestID] = rt
	}
}

func (ml *messengerLog) OnClose(topic string) {
	// pass
}

// getRequestRoundTrip gets the RoudnTrip object for a given requestID
func (ml *messengerLog) getRequestRoundTrip(requestID string) *RoundTrip {
	// If the requestID matches an open or closed request, return that. Otherwise return nil
	rt, ok := ml.openRequests[requestID]
	if ok {
		return rt
	}
	rt, ok = ml.completedRequests[requestID]
	if ok {
		return rt
	}
	return nil
}
