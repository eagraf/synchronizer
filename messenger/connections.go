package messenger

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// ConnectionManager maintains all connections for this server
// TODO determine whether this should be publicly accessible
type ConnectionManager struct {
	connections   map[string]*websocket.Conn
	subscriptions *PubSub
}

func newConnectionManager(ps *PubSub) *ConnectionManager {
	return nil
}

// AddConnection inserts the given connection, and begins listening for messages in a new goroutine
// TODO potentially handle promoting request as well
func (cm *ConnectionManager) AddConnection(workerUUID string, request *http.Request) {
	// Promote request to connection

	// Add to connections map

	// Begin listening

}

// RemoveConnection severs a websocket connection, and notifies all relevant listeners
func (cm *ConnectionManager) RemoveConnection(workerUUID string) {

}

// Send a message
func (cm *ConnectionManager) Send(workerUUID string, message *Message) {
	// Compress message
	// Send message
}

func (cm *ConnectionManager) listen() {
	// Wait for message on websocket

	// Deflate it

	// Notify relevant subscribers

}
