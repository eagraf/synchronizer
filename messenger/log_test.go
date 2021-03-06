package messenger

import "testing"

func TestStartRoundtrip(t *testing.T) {
	l := newMessengerLog()
	ps := newPubSub()

	ps.addTopic("new-topic")

	ps.addSubscription("new-topic", l)

	mb := MessageBuilder{}
	m, _ := mb.NewMessage("new-message", "request-id").Done()

	ps.publishSend("new-topic", m)
	rt, ok := l.openRequests["request-id"]
	if ok == false {
		t.Error("New roundtrip is missing")
	}
	if rt.Request == nil {
		t.Error("Request should be filled")
	}
	if rt.Response != nil {
		t.Error("Response should be nil")
	}

	if l.getRequestRoundTrip("request-id") == nil {
		t.Error("getRequestRoundtrip failed")
	}

	mb = MessageBuilder{}
	m, _ = mb.NewMessage("new-message", "request-id").Done()
	ps.publishReceived("new-topic", m)
	rt, ok = l.completedRequests["request-id"]
	if rt.Request == nil {
		t.Error("Request should be filled")
	}
	if rt.Response == nil {
		t.Error("Response should be filled")
	}

	if len(l.openRequests) != 0 {
		t.Error("Length of openRequests should be 0")
	}
	if len(l.completedRequests) != 1 {
		t.Error("Length of openRequests should be 0")
	}

	if l.getRequestRoundTrip("request-id") == nil {
		t.Error("getRequestRoundtrip failed")
	}
}

func TestEndRoundtrip(t *testing.T) {
	l := newMessengerLog()
	ps := newPubSub()

	ps.addTopic("new-topic")

	ps.addSubscription("new-topic", l)

	mb := MessageBuilder{}
	m, _ := mb.NewMessage("new-message", "request-id").Done()
	ps.publishReceived("new-topic", m)
	rt, ok := l.openRequests["request-id"]
	if rt.Request == nil {
		t.Error("Request should be filled")
	}
	if rt.Response != nil {
		t.Error("Response should be missing")
	}

	mb = MessageBuilder{}
	m, _ = mb.NewMessage("new-message", "request-id").Done()

	ps.publishSend("new-topic", m)
	rt, ok = l.completedRequests["request-id"]
	if ok == false {
		t.Error("New roundtrip is missing")
	}
	if rt.Request == nil {
		t.Error("Request should be filled")
	}
	if rt.Response == nil {
		t.Error("Response should be filled")
	}

	if len(l.openRequests) != 0 {
		t.Error("Length of openRequests should be 0")
	}
	if len(l.completedRequests) != 1 {
		t.Error("Length of openRequests should be 0")
	}
}
