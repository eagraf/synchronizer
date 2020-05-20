package selector

import (
	"fmt"

	"github.com/eagraf/synchronizer/messenger"
)

func (s *Selector) GetIdentifier() string {
	return "selector" // Temporary must fix TODO
}

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
func (s *Selector) OnSend(topic string, message *messenger.Message) {

}

func (s *Selector) OnClose(topic string) {
	fmt.Println(topic)
	s.workers[topic].Disconnected = true
	fmt.Println(s.workers[topic])
}
