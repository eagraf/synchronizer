package service

import (
	"net"
	"net/http"
	"strconv"

	"google.golang.org/grpc"
)

/*
 * The goal of this package is to provide a baseline Service object to different synchronizer
 * components. It helps with telemetry, health checks, and RPC service definitions.
 */

// Service represents a process running on a server
type Service struct {
	ID          string
	ServiceType string
	IP          *net.IP
	APIPort     int
	RPCPort     int
	Peers       map[string]map[string]*Service // Telemetry. Key1: Service Type, Key2: Service ID
	RPCService  interface{}
}

type Connection struct {
	service    *Service
	clientConn *grpc.ClientConn
}

// AddPeer registers a new peer for this service to communicate with
func (s *Service) AddPeer(newPeer *Service) error {
	return nil
}

// A ServiceInitiator is a driver for starting services (ServicePool or Kubernetes variant)
type ServiceInitiator interface {
	StartService(serviceType string, grpcServiceDescription *grpc.ServiceDesc, externalAPI *http.Handler) (*Service, error)
}

// ServicePool allows for multiple services to be interconnected and run locally without the need for starting a Kubernetes cluster
type ServicePool struct {
	// Implements ServiceInitiator
	// For development use only
	portCount int                            // Helper variable to keep track of unassigned ports
	Scale     scale                          // Number of each type of service
	Pool      map[string]map[string]*Service // Map of all services
}

type scale struct {
	test int
}

// NewServicePool creates a new ServicePool object
func NewServicePool() *ServicePool {
	sp := &ServicePool{
		portCount: 0,
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
		ID:          serviceType + " " + string(len(sp.Pool[serviceType])),
		ServiceType: serviceType,
		APIPort:     2000 + sp.portCount,
		RPCPort:     2001 + sp.portCount,
		Peers:       make(map[string]map[string]*Service),
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

// Get the RPC service description base off the service type
func getServiceDesc(serviceType string) *grpc.ServiceDesc {
	switch serviceType {
	case "Test":
		return &_Test_serviceDesc
	default:
		return nil
	}
}

// Helper for starting RPC server
func startRPCServer(service *Service, rpcHandler interface{}) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(service.RPCPort))
	if err != nil {
		return err
	}
	rpcServer := grpc.NewServer()
	rpcServer.RegisterService(getServiceDesc(service.ServiceType), rpcHandler)

	go rpcServer.Serve(listener)
	return nil
}

// Helper for starting external API
func startAPIServer(service *Service, apiHandler http.Handler) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(service.APIPort))
	if err != nil {
		return err
	}
	// TODO error handling if the server unexpectedly stops
	go http.Serve(listener, apiHandler)
	return nil
}
