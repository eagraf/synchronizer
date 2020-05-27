package service

import "testing"

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
