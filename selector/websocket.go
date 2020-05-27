package selector

import (
	"time"

	"github.com/eagraf/synchronizer/messenger"
	"github.com/google/uuid"
)

/*
 * websocket.go contains code for selector sending and receiving messages to/from workers
 */

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

// OnSend is a callback for outgoing messages
func (s *Selector) OnSend(topic string, message *messenger.Message) {

}

// OnClose is a callback for closed connections
func (s *Selector) OnClose(topic string) {
	s.workers[topic].Disconnected = true
}

/*
 * Types of sends:
 *   (1) Registration Response (selector_registration_response)
 *   (2) Health Check (selector_health_check)
 *   (3) Handoff (selector_handoff)
 */

func (s *Selector) sendRegistrationResponse() {

}

func (s *Selector) sendHealthCheck(workerUUID string) error {
	// Send the health check message
	mb := new(messenger.MessageBuilder)
	requestID := uuid.New().String()
	m, err := mb.NewMessage(MessageHealthCheck, requestID).Done()
	if err != nil {
		return err
	}
	s.messenger.Send(workerUUID, m)

	// Timeout waits in a new thread
	go func() {
		time.Sleep(HealthCheckTimeout)

		// Check if timeout was successful
		rt := s.messenger.GetRequestRoundTrip(requestID)
		if rt != nil && rt.Response != nil {
			s.workers[workerUUID].Healthy = true
		}
		s.workers[workerUUID].Healthy = false
	}()
	return nil
}

func (s *Selector) sendHandoff() {

}
