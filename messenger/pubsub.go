package messenger

// Subscriber is an interface implemented by routines that need to listen for specific messages
type Subscriber interface {
	GetIdentifier()                                          // Used to identify if the subscriber is already applied to a topic
	OnReceive(topic string, message *map[string]interface{}) // Callback after publish on message received
	OnSend(topic string, message *map[string]interface{})    // Callback after publish on message sent
	OnClose(topic string)                                    // Callback after topic closed
}

// PubSub handles notifying all relevant subroutines about incoming messages
type PubSub struct {
	subs map[string]map[string]*Subscriber
}

// AddSubscription applies a Subscriber to a topic
func (ps *PubSub) AddSubscription(topic string, subscriber *Subscriber) {
}

// RemoveSubscription removes a Subscriber from a topic
func (ps *PubSub) RemoveSubscription(topic string, subscriberID string) {

}

// TODO investigate if addTopic and closeTopic should be publicly available

// Create a new topic
func (ps *PubSub) addTopic(topic string) {

}

// Close a topic
func (ps *PubSub) closeTopic(topic string) {

}

// Notify all subscribers for a topic about some event
func (ps *PubSub) publish(topic string, message *Message) {

}
