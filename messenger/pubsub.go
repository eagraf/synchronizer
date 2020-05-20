package messenger

import (
	"errors"
)

// Subscriber is an interface implemented by routines that need to listen for specific messages
type Subscriber interface {
	GetIdentifier() string                    // Used to identify if the subscriber is already applied to a topic
	OnReceive(topic string, message *Message) // Callback after publish on message received
	OnSend(topic string, message *Message)    // Callback after publish on message sent
	OnClose(topic string)                     // Callback after topic closed
}

// pubSub handles notifying all relevant subroutines about incoming messages
type pubSub struct {
	subs map[string]map[string]*Subscriber
}

// Create a new pubsub
func newPubSub() *pubSub {
	ps := pubSub{
		subs: make(map[string]map[string]*Subscriber),
	}
	return &ps
}

// addSubscription applies a Subscriber to a topic
func (ps *pubSub) addSubscription(topic string, subscriber Subscriber) error {
	if _, ok := ps.subs[topic]; ok == false {
		return errors.New("Topic does not exist")
	}
	ps.subs[topic][subscriber.GetIdentifier()] = &subscriber
	return nil
}

// removeSubscription removes a Subscriber from a topic
func (ps *pubSub) removeSubscription(topic string, subscriberID string) error {
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
func (ps *pubSub) addTopic(topic string) {
	if _, ok := ps.subs[topic]; ok == false {
		ps.subs[topic] = make(map[string]*Subscriber)
	}
}

// Close a topic
func (ps *pubSub) closeTopic(topic string) {
	if _, ok := ps.subs[topic]; ok == true {
		for _, sub := range ps.subs[topic] {
			(*sub).OnClose(topic)
		}
		delete(ps.subs, topic)
	}
}

// Notify all relevant subscribers about a message received
func (ps *pubSub) publishReceived(topic string, message *Message) error {
	if _, ok := ps.subs[topic]; ok == false {
		return errors.New("Topic does not exist: " + topic)
	}

	for _, sub := range ps.subs[topic] {
		(*sub).OnReceive(topic, message)
	}
	return nil
}

// Notify all relevant subscribers about a message sent
func (ps *pubSub) publishSend(topic string, message *Message) error {
	if _, ok := ps.subs[topic]; ok == false {
		return errors.New("Topic does not exist: " + topic)
	}

	for _, sub := range ps.subs[topic] {
		(*sub).OnSend(topic, message)
	}
	return nil
}
