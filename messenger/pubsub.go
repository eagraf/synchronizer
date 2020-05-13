package messenger

import "errors"

// Subscriber is an interface implemented by routines that need to listen for specific messages
type Subscriber interface {
	GetIdentifier() string                    // Used to identify if the subscriber is already applied to a topic
	OnReceive(topic string, message *Message) // Callback after publish on message received
	OnSend(topic string, message *Message)    // Callback after publish on message sent
	OnClose(topic string)                     // Callback after topic closed
}

// PubSub handles notifying all relevant subroutines about incoming messages
type PubSub struct {
	subs map[string]map[string]*Subscriber
}

// Create a new pubsub
func newPubSub() *PubSub {
	ps := PubSub{
		subs: make(map[string]map[string]*Subscriber),
	}
	return &ps
}

// AddSubscription applies a Subscriber to a topic
func (ps *PubSub) AddSubscription(topic string, subscriber Subscriber) error {
	if _, ok := ps.subs[topic]; ok == false {
		return errors.New("Topic does not exist")
	}
	ps.subs[topic][subscriber.GetIdentifier()] = &subscriber
	return nil
}

// RemoveSubscription removes a Subscriber from a topic
func (ps *PubSub) RemoveSubscription(topic string, subscriberID string) error {
	if _, ok := ps.subs[topic]; ok == false {
		return errors.New("Topic does not exist")
	}

	if _, ok := ps.subs[topic][subscriberID]; ok == true {
		delete(ps.subs[topic], subscriberID)
	}
	return nil
}

// TODO investigate if addTopic and closeTopic should be publicly available

// Create a new topic
func (ps *PubSub) addTopic(topic string) {
	if _, ok := ps.subs[topic]; ok == false {
		ps.subs[topic] = make(map[string]*Subscriber)
	}
}

// Close a topic
func (ps *PubSub) closeTopic(topic string) {
	if _, ok := ps.subs[topic]; ok == true {
		for _, sub := range ps.subs[topic] {
			(*sub).OnClose(topic)
		}
		delete(ps.subs, topic)
	}
}

// Notify all relevant subscribers about a message received
func (ps *PubSub) publishReceived(topic string, message *Message) error {
	if _, ok := ps.subs[topic]; ok == false {
		return errors.New("Topic does not exist")
	}

	for _, sub := range ps.subs[topic] {
		(*sub).OnReceive(topic, message)
	}
	return nil
}

// Notify all relevant subscribers about a message sent
func (ps *PubSub) publishSend(topic string, message *Message) error {
	if _, ok := ps.subs[topic]; ok == false {
		return errors.New("Topic does not exist")
	}

	for _, sub := range ps.subs[topic] {
		(*sub).OnSend(topic, message)
	}
	return nil
}
