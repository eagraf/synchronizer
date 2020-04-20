package messenger

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gorilla/websocket"
)

type Messenger struct {
	connections   map[string]*websocket.Conn
	subscriptions map[string][]*Subscriber
}

type Message struct {
	MessageType string
	Payload     interface{}
}

type Subscriber interface {
	OnReceive(m *map[string]interface{})
}

var messenger *Messenger

// InitializeMessenger creates the messenger singleton
func InitializeMessenger() *Messenger {
	if messenger != nil {
		panic("Messenger already initialized")
	}
	m := Messenger{
		connections:   make(map[string]*websocket.Conn),
		subscriptions: make(map[string][]*Subscriber),
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
	m.connections[workerUUID] = connection
	go m.listen(workerUUID, connection)
}

// RemoveConnection stops listening on a connection
func (m *Messenger) RemoveConnection(workerUUID string) {
	delete(m.connections, workerUUID)
}

// AddSubscriber subscribes any subscriber interface to a topic
func (m *Messenger) AddSubscriber(subscriber *Subscriber, topics []string) {
	for _, topic := range topics {
		if _, ok := m.subscriptions[topic]; ok == false {
			m.subscriptions[topic] = make([]*Subscriber, 0)
		}
		m.subscriptions[topic] = append(m.subscriptions[topic], subscriber)
	}
}

// TODO RemoveSubscriber

// SendMessage sends a message to a worker in a separate thread
// TODO add some sort of error handling
func (m *Messenger) SendMessage(workerUUID string, payload interface{}) {
	go func() {
		// Construct the message
		message := Message{
			MessageType: reflect.TypeOf(payload).Elem().Name(),
			Payload:     payload,
		}
		// Encode message as json
		buffer, err := json.Marshal(message)
		if err != nil {
			fmt.Println(err.Error())
		}
		// Send the message over websocket
		err = m.connections[workerUUID].WriteMessage(websocket.TextMessage, buffer)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
}

// Listening thread for
func (m *Messenger) listen(workerUUID string, connection *websocket.Conn) {
	for {
		// Read incoming message
		_, buffer, err := connection.ReadMessage()
		// TODO check if error involves channel being closed
		if err != nil {
			fmt.Println(err)
			m.RemoveConnection(workerUUID)
			return
		}

		// Unmarshal the message
		var message interface{}
		err = json.Unmarshal(buffer, &message)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(buffer)
		fmt.Println(message)
		/*
			// Use workerUUID as a topic for now
			for _, subscriber := range m.subscriptions[workerUUID] {
				(*subscriber).OnReceive(message)
			}*/
	}
}
