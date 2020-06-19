package service

import (
	fmt "fmt"
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
	IP          net.IP
	APIPort     int
	RPCPort     int
	RPCService  interface{}
	peers       map[string]map[string]*Connection // Telemetry. Key1: Service Type, Key2: Service ID
}

// Connection represents a link between two services
type Connection struct {
	Service    *Service
	ClientConn *grpc.ClientConn
}

// AddPeer registers a new peer for this service to communicate with
func (s *Service) AddPeer(newPeer *Service) error {
	return nil
}

// A ServiceInitiator is a driver for starting services (ServicePool or Kubernetes variant)
type ServiceInitiator interface {
	StartService(serviceType string, rpcHandler interface{}, apiHandler http.Handler) (*Service, error)
	ConnectService(service *Service) error
}

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

// DefaultTopology is a base topology that can be used for testing
var DefaultTopology map[string]map[string]bool = map[string]map[string]bool{
	"Test": {
		"Test": true,
	},
	"Coordinator": {
		"Selector": true,
	},
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
	fmt.Println(service.ID)
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

// Establish a RPC connection
// TODO handle unexpected connection failures
func connect(source *Service, dest *Service) error {
	conn, err := grpc.Dial(dest.IP.String()+":"+strconv.Itoa(dest.RPCPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	// Add to peers
	c := &Connection{
		Service:    dest,
		ClientConn: conn,
	}

	// Create new map if needed for the service type
	if _, ok := source.peers[dest.ServiceType]; ok == false {
		source.peers[dest.ServiceType] = make(map[string]*Connection)
	}

	source.peers[dest.ServiceType][dest.ID] = c
	return nil
}

// Get the RPC service description base off the service type
func getServiceDesc(serviceType string) *grpc.ServiceDesc {
	switch serviceType {
	case "Test":
		return &_Test_serviceDesc
	case "Selector":
		return &_Selector_serviceDesc
	case "Data_Server":
		return &_DataServerService_serviceDesc
	case "Aggregator":
		return &_Aggregator_serviceDesc
	default:
		return nil
	}
}

// Helper for starting RPC server
func startRPCServer(service *Service, rpcHandler interface{}) error {
	// If rpcHandler is nil, do nothing
	if rpcHandler == nil {
		return nil
	}
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
