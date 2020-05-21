package selector

import (
	"github.com/eagraf/synchronizer/messenger"
)

// GetIdentifier returns the selector's id
func (s *Selector) GetIdentifier() string {
	return "selector" // Temporary must fix TODO
}

// OnReceive processes incoming messages
func (s *Selector) OnReceive(topic string, message *messenger.Message) {
	switch message.GetMetadata().MessageType {
	case MessageHealthCheck:
	}
}

/*
 * Types of sends:
 *   (1) Registration Response (selector_registration_response)
 *   (2) Health Check (selector_health_check)
 *   (3) Handoff (selector_handoff)
 */

// OnSend is a callback for outgoing messages
func (s *Selector) OnSend(topic string, message *messenger.Message) {

}

// OnClose is a callback for closed connections
func (s *Selector) OnClose(topic string) {
	s.workers[topic].Disconnected = true
}
