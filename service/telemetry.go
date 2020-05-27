package service

import "errors"

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
