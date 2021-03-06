package messenger

import "net/http"

/*
 * Reusable package for maintaining websocket connections in go services
 * Features:
 * 	(1) Extendible message type with a header for structured data, and payload section for raw data that needs quick copying
 *  (2) Pubsub system that allows for different modules to subscribe to connections and listen for messages
 *  (3) Connection management with message loss detection (TODO)
 *  (4) Message and request logging (TODO)
 */

/*
 * Exported types from this package
 * 	(1) Messenger
 * 	(2) Subscriber struct
 * 	(3) Message and related types
 * 	(4) MessageBuilder
 *  (5) TestClient
 *  (6) RoundTrip
 */

// Messenger type is a wrapper for ConnectionManager, PubSub, and more
type Messenger struct {
	cm *connectionManager
	ps *pubSub
	ml *messengerLog
}

// NewMessenger initializes a new messenger instance
func NewMessenger() *Messenger {
	ps := newPubSub()
	cm := newConnectionManager(ps)
	ml := newMessengerLog()
	m := Messenger{
		cm,
		ps,
		ml,
	}
	return &m
}

// connectionManager wrappers

// AddConnection is a wrapper method for connectionManager.AddConnection
func (m *Messenger) AddConnection(workerUUID string, writer http.ResponseWriter, request *http.Request) error {
	err := m.cm.AddConnection(workerUUID, writer, request)
	// Have messenger logger listen to updates on this connection as well
	if err != nil {
		m.AddSubscription(workerUUID, m.ml)
	}
	return err
}

// RemoveConnection is a wrapper method for connectionManager.RemoveConnection
func (m *Messenger) RemoveConnection(workerUUID string) {
	m.cm.RemoveConnection(workerUUID)
}

// Send is a wrapper method for connectionManager.Send
func (m *Messenger) Send(workerUUID string, message *Message) {
	m.cm.Send(workerUUID, message)
}

// pubSub wrappers

// AddSubscription is a wrapper method for pubSub.AddSubscription
func (m *Messenger) AddSubscription(topic string, subscriber Subscriber) error {
	return m.ps.addSubscription(topic, subscriber)
}

// RemoveSubscription is a wrapper method for pubSub.removeSubscription
func (m *Messenger) RemoveSubscription(topic string, subscriberID string) error {
	return m.ps.removeSubscription(topic, subscriberID)
}

// messengerLog wrappers

// GetRequestRoundTrip wraps messengerLog.getRequestRoundTrip
func (m *Messenger) GetRequestRoundTrip(requestID string) *RoundTrip {
	return m.ml.getRequestRoundTrip(requestID)
}