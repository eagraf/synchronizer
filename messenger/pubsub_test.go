package messenger

import (
	"testing"
)

type subscriberImpl struct {
	id string
}

func (si *subscriberImpl) GetIdentifier() string {
	return si.id
}

func (si *subscriberImpl) OnReceive(topic string, message *Message) {
	onReceives++
}

func (si *subscriberImpl) OnSend(topic string, message *Message) {
	onSends++
}

func (si *subscriberImpl) OnClose(topic string) {
	onCloses++
}

// Globals tickers used to test callback invocations
var (
	onReceives int = 0
	onSends    int = 0
	onCloses   int = 0
)

func TestAddCloseTopic(t *testing.T) {
	ps := newPubSub()

	ps.addTopic("new-topic")
	if len(ps.subs) != 1 {
		t.Error("Wrong number of topics")
	}

	ps.addTopic("second-topic")
	if len(ps.subs) != 2 {
		t.Error("Wrong number of topics")
	}

	// Try adding existing topic
	ps.addTopic("new-topic")
	if len(ps.subs) != 2 {
		t.Error("Wrong number of topics")
	}

	ps.closeTopic("new-topic")
	if len(ps.subs) != 1 {
		t.Error("Wrong number of topics")
	}
}

func TestCloseTopicMissing(t *testing.T) {
	// Defined behavior: if topic is missing, don't error
	ps := newPubSub()

	ps.addTopic("new-topic")
	ps.closeTopic("second-topic")
	if len(ps.subs) != 1 {
		t.Error("Wrong number of topics")
	}
}

func TestAddRemoveSubscription(t *testing.T) {
	ps := newPubSub()

	ps.addTopic("new-topic")

	s1 := subscriberImpl{"a"}
	s2 := subscriberImpl{"b"}
	ps.addSubscription("new-topic", &s1)
	if len(ps.subs["new-topic"]) != 1 {
		t.Error("Wrong number of subscribers")
	}

	ps.addSubscription("new-topic", &s2)
	if len(ps.subs["new-topic"]) != 2 {
		t.Error("Wrong number of subscribers")
	}

	ps.removeSubscription("new-topic", "a")
	if len(ps.subs["new-topic"]) != 1 {
		t.Error("Wrong number of subscribers")
	}
	ps.removeSubscription("new-topic", "b")
	if len(ps.subs["new-topic"]) != 0 {
		t.Error("Wrong number of subscribers")
	}
}

func TestAddSubMissingTopic(t *testing.T) {
	ps := newPubSub()

	s1 := subscriberImpl{"a"}
	err := ps.addSubscription("new-topic", &s1)
	if err == nil {
		t.Error("Error expected")
	}
}

func TestMultiTopicAddRemoveSubscription(t *testing.T) {
	ps := newPubSub()

	ps.addTopic("topic1")
	ps.addTopic("topic2")

	s1 := subscriberImpl{"a"}
	ps.addSubscription("topic1", &s1)
	ps.addSubscription("topic2", &s1)
	if len(ps.subs["topic1"]) != 1 {
		t.Error("Wrong number of subscribers")
	}
	if len(ps.subs["topic2"]) != 1 {
		t.Error("Wrong number of subscribers")
	}

	ps.removeSubscription("topic1", "a")
	if len(ps.subs["topic1"]) != 0 {
		t.Error("Wrong number of subscribers")
	}
	if len(ps.subs["topic2"]) != 1 {
		t.Error("Wrong number of subscribers")
	}
}

func TestPublishReceive(t *testing.T) {
	ps := newPubSub()

	ps.addTopic("new-topic")

	s1 := subscriberImpl{"a"}
	s2 := subscriberImpl{"b"}
	ps.addSubscription("new-topic", &s1)
	ps.addSubscription("new-topic", &s2)

	mb := MessageBuilder{}
	m, _ := mb.NewMessage("new-message", "request-id").Done()

	ps.publishReceived("new-topic", m)
	if onReceives != 2 {
		t.Error("Incorrect number of onReceive invocations")
	}
}

func TestPublishSend(t *testing.T) {
	ps := newPubSub()

	ps.addTopic("new-topic")

	s1 := subscriberImpl{"a"}
	s2 := subscriberImpl{"b"}
	ps.addSubscription("new-topic", &s1)
	ps.addSubscription("new-topic", &s2)

	mb := MessageBuilder{}
	m, _ := mb.NewMessage("new-message", "request-id").Done()

	ps.publishSend("new-topic", m)
	if onSends != 2 {
		t.Error("Incorrect number of onReceive invocations")
	}
}

func TestCloseTopicSubscribers(t *testing.T) {
	ps := newPubSub()

	ps.addTopic("new-topic")

	s1 := subscriberImpl{"a"}
	s2 := subscriberImpl{"b"}
	ps.addSubscription("new-topic", &s1)
	ps.addSubscription("new-topic", &s2)

	ps.closeTopic("new-topic")
	if onCloses != 2 {
		t.Error("Incorrect number of onClose invocations")
	}
}
