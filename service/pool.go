package service

import (
	"net"
	"net/http"
	"strconv"
)

// This file contains an implementation of ServiceInitiator used for local testing

// ServicePool allows for multiple services to be interconnected and run locally without the need for starting a Kubernetes cluster
type ServicePool struct {
	// Implements ServiceInitiator
	// For development use only
	portCount   int // Helper variable to keep track of unassigned ports
	initialPort int
	topology    map[string]map[string]bool
	//	Scale     scale                          // Number of each type of service
	Pool map[string]map[string]*Service // Map of all services
}

// NewServicePool creates a new ServicePool object
func NewServicePool(initialPort int, top map[string]map[string]bool) *ServicePool {
	sp := &ServicePool{
		portCount:   0,
		initialPort: initialPort,
		topology:    top, // Key 1: Origin Service, Key 2: Receiving Service, Value: There is a link
		//Scale:     scale{},
		Pool: make(map[string]map[string]*Service),
	}
	return sp
}

/*
 * Stages of starting service:
 *   (1) Start service servers (external API and RPC)
 *   (2) Connect to other services
 */

// StartService creates a new service and connects it to the correct peer services
func (sp *ServicePool) StartService(serviceType string, rpcHandler interface{}, apiHandler http.Handler) (*Service, error) {
	// Check service is valid
	if _, ok := sp.Pool[serviceType]; ok == false {
		sp.Pool[serviceType] = make(map[string]*Service)
	}

	// Create Service
	service := &Service{
		ID:          serviceType + "-" + strconv.Itoa(len(sp.Pool[serviceType])), // TODO this needs to be based off of an incrementing counter (currently breaks when the service count goes down)
		IP:          net.IPv4(127, 0, 0, 1),
		ServiceType: serviceType,
		APIPort:     sp.initialPort + sp.portCount,
		RPCPort:     sp.initialPort + sp.portCount + 1,
		peers:       make(map[string]map[string]*Connection),
	}

	// Update portCount
	sp.portCount += 2
	// Update the service pool
	sp.Pool[serviceType][service.ID] = service

	// Start servers
	if err := startAPIServer(service, apiHandler); err != nil {
		return nil, err
	}

	if err := startRPCServer(service, rpcHandler); err != nil {
		return nil, err
	}

	return service, nil
}

// ConnectService establishes RPC connections based off of the service topology
func (sp *ServicePool) ConnectService(service *Service) error {
	// TODO handle errors
	// Connect incoming services
	for st, connections := range sp.topology { // Range through each type of service
		if _, ok := connections[service.ServiceType]; ok == true {
			// Establish connections
			for _, s := range sp.Pool[st] { // Range through all services of a type
				connect(s, service)
			}
		}
	}

	// Connect outgoing services
	for st := range sp.topology[service.ServiceType] { // Range through each outgoing service connection
		for _, s := range sp.Pool[st] { // Range through all services of a type
			connect(service, s)
		}
	}
	return nil
}
