package messenger

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// ConnectionManager maintains all connections for this server
// TODO determine whether this should be publicly accessible
type ConnectionManager struct {
	connections   map[string]*connection
	subscriptions *PubSub
}

type connection struct {
	connection *websocket.Conn
	mutex      sync.Mutex
}

// Create new connectionManager
func newConnectionManager(ps *PubSub) *ConnectionManager {
	cm := ConnectionManager{
		connections:   make(map[string]*connection),
		subscriptions: ps,
	}
	return &cm
}

// AddConnection inserts the given connection, and begins listening for messages in a new goroutine
func (cm *ConnectionManager) AddConnection(workerUUID string, writer http.ResponseWriter, request *http.Request) error {
	// Promote request to connection
	// upgrader promotes a standard HTTP/HTTPS connection to a websocket connection
	// TODO implement CheckOrigin
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return err
	}

	// Add to connections map
	if _, ok := cm.connections[workerUUID]; ok == true {
		return errors.New("Connection with this UUID already exists")
	}
	cm.connections[workerUUID] = &connection{
		connection: conn,
		mutex:      sync.Mutex{},
	}

	// Add topic that can be subscribed too
	cm.subscriptions.addTopic(workerUUID)

	// Begin listening
	go cm.listen(workerUUID, cm.connections[workerUUID])

	return nil
}

// RemoveConnection severs a websocket connection, and notifies all relevant listeners
func (cm *ConnectionManager) RemoveConnection(workerUUID string) {
	if _, ok := cm.connections[workerUUID]; ok == true {
		// Must be able to acquire a lock before deleting
		// TODO check this doesn't cause a race condition if message was sent right beforehand. This could occur if there is a context switch from Send to RemoveConnection before Send acquires lock
		c := cm.connections[workerUUID]
		c.mutex.Lock()
		delete(cm.connections, workerUUID)
		c.connection.Close() // Close websocket, no need to check for error
		c.mutex.Unlock()
	}

	cm.subscriptions.closeTopic(workerUUID)
}

// Send a message
// This method cannot return an error without causing a wait.
// TODO error logging
// TODO recovery from failed sends
func (cm *ConnectionManager) Send(workerUUID string, message *Message) {
	if _, ok := cm.connections[workerUUID]; ok == false {
		fmt.Println(errors.New("No connection with that UUID exists"))
	}

	// Send message in new thread
	go func() {
		// Lock needs to be established before calling NextWriter
		cm.connections[workerUUID].mutex.Lock()
		// Get the writer, to be passed into the compression writer
		w, err := cm.connections[workerUUID].connection.NextWriter(websocket.BinaryMessage)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Send the message
		// Compresses at send time
		err = writeMessage(message, w)
		// Closing writer tells websocket to send message
		w.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		cm.connections[workerUUID].mutex.Unlock()

		// Notify subscribers of send
		err = cm.subscriptions.publishSend(workerUUID, message)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
}

// Main listening loop
func (cm *ConnectionManager) listen(workerUUID string, c *connection) {

	for {
		// Compressed message is read into buffer from websocket
		_, buffer, err := c.connection.ReadMessage()
		if err != nil {
			// Prevent concurrent writes
			c.mutex.Lock()
			defer c.mutex.Unlock()

			// If the websocket was closed by the connectionManager, subscribers have already been notified
			if _, ok := cm.connections[workerUUID]; ok == false {
				return
			}
			// Otherwise handle error
			// TODO explicitly handle all possible errors
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseAbnormalClosure:
					cm.subscriptions.closeTopic(workerUUID)
					fmt.Println(err.Error())
					return
				default:
					cm.subscriptions.closeTopic(workerUUID)
					fmt.Println(err.Error())
					return

				}
			} else {
				cm.subscriptions.closeTopic(workerUUID)
				fmt.Println(err.Error())
				return
			}
		}
		// Message is inflated
		message, err := readMessage(buffer)
		if err != nil {
			fmt.Println(err.Error())
		}

		// Notifying all subscribers
		err = cm.subscriptions.publishReceived(workerUUID, message)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
