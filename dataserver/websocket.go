package dataserver

import "github.com/eagraf/synchronizer/messenger"

/*
 * websocket.go contains code for selector sending and receiving messages to/from workers
 */

// GetIdentifier returns the selector's id
func (ds *DataServer) GetIdentifier() string {
	return "data_server" // Temporary must fix TODO
}

// OnReceive processes incoming messages
func (ds *DataServer) OnReceive(topic string, message *messenger.Message) {
}

// OnSend is a callback for outgoing messages
func (ds *DataServer) OnSend(topic string, message *messenger.Message) {
}

// OnClose is a callback for closed connections
func (ds *DataServer) OnClose(topic string) {
}
