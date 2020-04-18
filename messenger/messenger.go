package messenger

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gorilla/websocket"
)

type Messenger struct {
	connections map[string]*websocket.Conn
}

type Message struct {
	MessageType string
	Payload     interface{}
}

var messenger *Messenger

// InitializeMessenger creates the messenger singleton
func InitializeMessenger() *Messenger {
	if messenger != nil {
		panic("Messenger already initialized")
	}
	m := Messenger{
		connections: make(map[string]*websocket.Conn),
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
	fmt.Println("Listening")
	for {
		_, buffer, err := connection.ReadMessage()
		// TODO check if error involves channel being closed
		if err != nil {
			fmt.Println(err)
			m.RemoveConnection(workerUUID)
			return
		}

		fmt.Println(buffer)

		var message interface{}
		err = json.Unmarshal(buffer, &message)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(message)
	}
}