package aggregator

import "github.com/eagraf/synchronizer/messenger"

/*
 * websocket.go contains code for selector sending and receiving messages to/from workers
 */

// GetIdentifier returns the selector's id
func (ds *Aggregator) GetIdentifier() string {
	return "aggregator" // Temporary must fix TODO
}

// OnReceive processes incoming messages
func (ds *Aggregator) OnReceive(topic string, message *messenger.Message) {
}

// OnSend is a callback for outgoing messages
func (ds *Aggregator) OnSend(topic string, message *messenger.Message) {
}

// OnClose is a callback for closed connections
func (ds *Aggregator) OnClose(topic string) {
}
