package messenger

import "time"

// TODO currently the roundtrip request/response may not resemble all patterns of communication. A better system for interlinking messages may be needed (graph?)
type log struct {
	openRequests  map[string]*roundTrip
	requestBuffer []*RoundTrip // TODO this will grow out of control -> a better logging system is needed
	logger        log.Logger
}

// RoundTrip represents a request response pair
type RoundTrip struct {
	requestID string
	request   *Message
	response  *Message
}

// Duration gets the total duration of a request
func (rt *RoundTrip) Duration() time.Duration {
	return time.Parse(time.StampMilli, rt.response.metadata.Timestamp).
		Sub(time.Parse(time.StampMilli, rt.request.metadata.Request.metadata.Timestamp)).
		Format(time.StampMilli)
}

func (l *log) GetIdentifier() string {
	// TODO universal identifier
	return "messenger_logger"
}

// OnReceive callback insert or complete roundtrip
func (l *log) OnReceive(topic string, message *Message) {
	// Check if it matches any active message request
	rt, ok := l.openRequests[topic]
	if ok == true {
		// Complete round trip
		rt.response = message
	} else {
		// Otherwise create new request
		rt = new(RoundTrip)
		rt.requestID = message.metadata.Request
		l.openRequests[rt.requestID] = rt
	}

}

// OnSend callback insert or complete roundtrip
func (l *Log) OnSend(topic string, message *Message) {
	// Check if it matches any active message request
	rt, ok := l.openRequests[topic]
	if ok == true {
		// Complete round trip
		rt.response = message
	} else {
		// Otherwise create new request
		rt = new(RoundTrip)
		rt.requestID = message.metadata.Request
		l.openRequests[rt.requestID] = rt
	}
}

func (l *Log) OnClose(topic string) {
	// pass
}
