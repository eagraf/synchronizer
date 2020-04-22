package messenger

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Messenger struct {
	connections    map[string]*Connection
	subscriptions  map[string]map[string]*Subscriber
	activeRequests map[string]Request
}

type Message struct {
	MessageType string
	Payload     interface{}
}

type Connection struct {
	connection *websocket.Conn
	mutex      sync.Mutex
}

type Subscriber interface {
	AddSubscription(topic string)
	GetUUID() string
	OnReceive(topic string, m *map[string]interface{})
	OnClose(string)
}

type Request struct {
	start int64
	end   int64
}

var messenger *Messenger

// InitializeMessenger creates the messenger singleton
func InitializeMessenger() *Messenger {
	if messenger != nil {
		panic("Messenger already initialized")
	}
	m := Messenger{
		connections:    make(map[string]*Connection),
		subscriptions:  make(map[string]map[string]*Subscriber), // Each topic is a set of subscribers
		activeRequests: make(map[string]Request),
	}
	messenger = &m
	return messenger
}

// GetMessengerSingleton returns the messenger
func GetMessengerSingleton() *Messenger {
	if messenger == nil {
		panic("Messenger has not been initialized yet")
	}
	return messenger
}

// AddConnection registers a connection with the messenger and begins listening on a new thread
func (m *Messenger) AddConnection(workerUUID string, connection *websocket.Conn) {
	c := Connection{
		connection,
		sync.Mutex{},
	}
	m.connections[workerUUID] = &c
	go m.listen(workerUUID, &c)
}

// RemoveConnection stops listening on a connection
func (m *Messenger) RemoveConnection(workerUUID string) {
	// Delete connection
	delete(m.connections, workerUUID)

	// Notify any listeners
	if _, ok := m.subscriptions[workerUUID]; ok == true {
		for i, subscriber := range m.subscriptions[workerUUID] {
			(*subscriber).OnClose(workerUUID)
		}
	}

	// Delete topic
	delete(m.subscriptions, workerUUID)
}

// AddSubscriber subscribes any subscriber interface to a topic
func (m *Messenger) AddSubscriber(subscriber Subscriber, topics []string) {
	for _, topic := range topics {
		// Check if topic exists
		if _, ok := m.subscriptions[topic]; ok == false {
			// If not make topic
			m.subscriptions[topic] = make(map[string]*Subscriber)
		}
		// Check if already subscribed
		if _, ok := m.subscriptions[topic][subscriber.GetUUID()]; ok == false {
			m.subscriptions[topic][subscriber.GetUUID()] = &subscriber
			subscriber.AddSubscription(topic)
		}
	}
}

// RemoveSubscriber unsubscribes a listener
func (m *Messenger) RemoveSubscriber(subscriber Subscriber, topic string) {
	if _, ok := m.subscriptions[topic][subscriber.GetUUID()]; ok == true {
		delete(m.subscriptions[topic], subscriber.GetUUID())
	}
}

// SendMessage sends a message to a worker in a separate thread
// TODO add some sort of error handling
// TODO account for concurrent write to websocket connection error
func (m *Messenger) SendMessage(workerUUID string, payload interface{}) {
	go func() {
		messageType := reflect.TypeOf(payload).Elem().Name()
		// Construct the message
		message := Message{
			MessageType: messageType,
			Payload:     payload,
		}
		// Encode message as json
		buffer, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Add to request queue
		if messageType == "Intent" {
			r := Request{
				start: time.Now().UnixNano() / int64(time.Millisecond),
			}
			m.activeRequests[workerUUID] = r
		}

		// Send the message over websocket
		m.connections[workerUUID].mutex.Lock()
		err = m.connections[workerUUID].connection.WriteMessage(websocket.TextMessage, buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
		m.connections[workerUUID].mutex.Unlock()
	}()
}

// Listening thread for
func (m *Messenger) listen(workerUUID string, c *Connection) {
	for {
		// Read incoming message
		_, buffer, err := c.connection.ReadMessage()
		// TODO check if error involves channel being closed
		if err != nil {
			fmt.Println("Removing connection:" + err.Error())
			m.RemoveConnection(workerUUID)
			return
		}

		// Unmarshal the message
		var message map[string]interface{}
		err = json.Unmarshal(buffer, &message)
		if err != nil {
			fmt.Println(err)
		}

		if s, ok := m.activeRequests[workerUUID]; ok != false {
			s.end = time.Now().UnixNano() / int64(time.Millisecond)
			fmt.Println(s)
		}
		// Why does everything suck
		delete(m.activeRequests, workerUUID)
		// Use workerUUID as a topic for now
		for _, subscriber := range m.subscriptions[workerUUID] {
			(*subscriber).OnReceive(workerUUID, &message)
		}
	}
}
