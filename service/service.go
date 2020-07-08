package service

import (
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
	IP          net.IP
	APIPort     int
	RPCPort     int
	RPCService  interface{}
	Logger      *log.Logger
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

// DefaultTopology is a base topology that can be used for testing
var DefaultTopology map[string]map[string]bool = map[string]map[string]bool{
	"Test": {
		"Test": true,
	},
	"Coordinator": {
		"Selector": true,
	},
}

// Below are a bunch of helper functions (TODO need to be moved for better code organization)

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
