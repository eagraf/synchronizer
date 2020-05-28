package service

import (
	context "context"
	"errors"
	"time"
)

/*
 * Provide helper methods for services to identify the correct peers to communicate with
 * when state is involved.
 */

// ConnectionSet represents a set of active connections that the service is maintaining
type ConnectionSet map[string]*Connection

// AllPeersOfType returns a map of all connected services of a given type
func (s *Service) AllPeersOfType(serviceType string) (ConnectionSet, error) {
	if ps, ok := s.peers[serviceType]; ok == true {
		return ps, nil
	}

	return nil, errors.New("No peers of the type " + serviceType)
}

// Make a request in a new thread using invoke
func (c *Connection) RPCRequest(method string, args, reply interface{}) chan interface{} {

	replyChannel := make(chan interface{})
	go func() {
		cc := c.ClientConn
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := cc.Invoke(ctx, method, args, reply)
		if err != nil {
			replyChannel <- err
		}
		replyChannel <- reply
	}()

	return replyChannel
}

// This is a blocking call
// Each service should be of the same type
/*func (s *Service) MultiCast(targets ConnectionSet, method string, args interface{}, replys []interface{}) ([]interface{}, error) {
	// Timeout???

	responseChans := make([]reflect.SelectCase, len(targets))
	errChans := make([]reflect.SelectCase, len(targets))
	i := 0
	for _, t := range targets {
		responseChan, errChan := t.RPCRequest(method, args, replys[i])
		i++
	}

	responses := make([]reflect.SelectCase, len(targets))
	for _, t := range targets {
		responses[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	chosen, value, ok := reflect.Select(cases)
	// ok will be true if the channel has not been closed.
	ch := chans[chosen]
	msg := value.String()
}*/
