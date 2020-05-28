package service

import (
	"reflect"
	"testing"
)

func TestAllPeersOfType(t *testing.T) {
	// Dummy service
	s := &Service{
		peers: make(map[string]map[string]*Connection),
	}
	s.peers["service1"] = map[string]*Connection{
		"peer1": {
			nil,
			nil,
		},
		"peer2": {
			nil,
			nil,
		},
	}
	ps, err := s.AllPeersOfType("service1")
	if err != nil {
		t.Error(err.Error())
	}
	if len(ps) != 2 {
		t.Error("Incorrect number of peers")
	}

	ps, err = s.AllPeersOfType("service2")
	if err == nil {
		t.Error("Expected nil error")
	}
}

func TestRPCRequest(t *testing.T) {
	topology := make(map[string]map[string]bool)
	topology["Test"] = map[string]bool{"Test": true}
	sp = NewServicePool(2100, topology)

	// Create new TestService
	t1, err := createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	t2, err := createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	sp.ConnectService(t2)

	reply := Pong{}
	replyChan := t2.peers["Test"][t1.ID].RPCRequest("/testservice.Test/TestRPC", &Ping{Message: "Hello"}, &reply)

	select {
	case r := <-replyChan:
		if reflect.TypeOf(r).String() == "*status.Error" {
			t.Error(r.(error).Error())
		} else {
			if r.(*Pong).Message != "Hello" {
				t.Error("Incorrect message value")
			}
		}
	}

	// Try making incorrect request
	replyChan = t2.peers["Test"][t1.ID].RPCRequest("bad_method", &Ping{Message: "Hello"}, &reply)

	select {
	case r := <-replyChan:
		if reflect.TypeOf(r).String() != "*status.Error" {
			t.Error("Expected an error")
		}
	}

}
