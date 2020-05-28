package service

import (
	context "context"
	"errors"
	fmt "fmt"
	"reflect"
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

// RPCRequest invokes a gRPC request in a new thread, and returns the result through a channel
// Response interace needs to be checked as either the reply value or as an error
// Make a request in a new thread using invoke
func (c *Connection) RPCRequest(method string, args, reply interface{}) chan interface{} {

	replyChannel := make(chan interface{})
	// Send message in new thread
	go func() {
		cc := c.ClientConn
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Invoke RPC method
		err := cc.Invoke(ctx, method, args, reply)
		// Send response via channel
		if err != nil {
			replyChannel <- err
		} else {
			replyChannel <- reply
		}
	}()

	return replyChannel
}

// MultiCast sends an identical gRPC request to multiple services concurrently, blocking until all have returned
// This is a blocking call
// Each service should be of the same type
// If there are errors, the corresponding index in replys will be set to the error value, and the error will be appended to the error slice returned
func (s *Service) MultiCast(targets ConnectionSet, method string, args interface{}, replys []interface{}) ([]interface{}, []error) {
	// TODO add timeout

	// Make requests to each service
	// Uses reflec.SelectCase to handle dynamic selects
	responseChans := make([]reflect.SelectCase, len(targets))
	i := 0
	for _, t := range targets {
		responseChan := t.RPCRequest(method, args, replys[i])
		responseChans[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(responseChan)}
		i++
	}

	// This loop blocks while all responses are being received
	errs := make([]error, 0, len(targets))
	responseCount := 0
	for range targets {
		responseCount++
		// TODO handle case if channel has been closed
		i, value, _ := reflect.Select(responseChans)

		// If it is an error, handle appropriately
		if _, isErr := value.Interface().(error); isErr == true {
			errs = append(errs, fmt.Errorf("Index %d: %s", i, value.Interface().(error).Error()))
		}
		replys[i] = value.Interface()
	}
	return replys, errs
}
