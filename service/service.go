package service

import (
	fmt "fmt"
	"log"
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

type ServiceInitiator interface {
	StartService(serviceType string, grpcServiceDescription *grpc.ServiceDesc, externalAPI *http.Handler) (*Service, error)
}

// ServicePool allows for multiple services to be interconnected and run locally without the need for starting a Kubernetes cluster
// Implements ServiceInitiator
// For development use only
type ServicePool struct {
	portCount int // Helper variable to keep track of unassigned ports
	log       log.Logger
	Scale     scale                          // Number of each type of service
	Pool      map[string]map[string]*Service // Map of all services
}

type scale struct {
	test int
}

func NewServicePool() *ServicePool {
	sp := &ServicePool{
		portCount: 0,
		//Scale:     scale{},
		Pool: make(map[string]map[string]*Service),
	}
	return sp
}

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

	// Start external API handling
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(service.APIPort))
	if err != nil {
		log.Fatal("listen error (api):", err)
		return nil, err
	}
	go http.Serve(listener, apiHandler)

	//go log.Fatal(http.ListenAndServe(":"+strconv.Itoa(service.APIPort), apiHandler))
	fmt.Println("yo")

	// Start handling gRPC
	go func() {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(service.RPCPort))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		rpcServer := grpc.NewServer()
		rpcServer.RegisterService(getServiceDesc(serviceType), rpcHandler)

		//pb.RegisterGreeterServer(s, &server{})
		if err := rpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

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
