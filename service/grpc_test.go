package service

import (
	"reflect"
	"testing"
)

func TestGetPeer(t *testing.T) {
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

	_, err := s.GetPeer("service1", "peer1")
	if err != nil {
		t.Error("Did not expect an error")
	}

	_, err = s.GetPeer("blah", "peer1")
	if err == nil {
		t.Error("Expected an error")
	}

	_, err = s.GetPeer("service1", "blah")
	if err == nil {
		t.Error("Expected an error")
	}
}

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
	sp = NewServicePool(3000, topology)

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

func TestUniCast(t *testing.T) {
	topology := make(map[string]map[string]bool)
	topology["Test"] = map[string]bool{"Test": true}
	sp = NewServicePool(3010, topology)

	// Create new TestService
	_, err := createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	t2, err := createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	sp.ConnectService(t2)

	closureTest := 0

	reply := Pong{}
	conn := t2.peers["Test"]["Test-1"]
	t2.UniCast(conn, "/testservice.Test/TestRPC", &Ping{Message: "Hello"}, &reply, func(reply interface{}, err error) {
		closureTest++
		if closureTest != 1 {
			t.Error("Closure test failed")
		}
		if reply.(*Pong).Message != "Hello" {
			t.Error("Incorrect response messsage value")
		}
	})

	// Test bad method
	reply = Pong{}
	t2.UniCast(conn, "bad method", &Ping{Message: "Hello"}, &reply, func(reply interface{}, err error) {
		if reply != nil {
			t.Error("reply should be nil")
		}
		if err == nil {
			t.Error("err should not be nil")
		}
	})

}

func TestMultiCast(t *testing.T) {
	topology := make(map[string]map[string]bool)
	topology["Test"] = map[string]bool{"Test": true}
	sp = NewServicePool(3020, topology)

	// Create new TestService
	_, err := createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	_, err = createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	t4, err := createTestServerImpl(sp)
	if err != nil {
		t.Error(err.Error())
	}
	sp.ConnectService(t4)

	replys := make([]interface{}, 4)
	for i := range replys {
		replys[i] = &Pong{}
	}
	peers, _ := t4.AllPeersOfType("Test")
	responses, errs := t4.MultiCast(peers, "/testservice.Test/TestRPC", &Ping{Message: "Hello"}, replys)
	if len(errs) != 0 {
		t.Error(errs[0].Error())
	}
	if len(responses) != 4 {
		t.Error("Incorrect number of responses")
	}
	if responses[0].(*Pong).Message != "Hello" {
		t.Error("Incorrect response messsage value")
	}

	for i := range replys {
		replys[i] = &Pong{}
	}
	// Test errors
	responses, errs = t4.MultiCast(peers, "bad_method", &Ping{Message: "Hello"}, replys)
	if len(errs) != 4 {
		t.Errorf("Incorrect number of errors: %d", len(errs))
	}
	if len(responses) != 4 {
		t.Error("Incorrect responses length")
	}
	if _, isErr := responses[0].(error); isErr == false {
		t.Errorf("Response is not of type error: %v", reflect.TypeOf(responses[0]))
	}
}
