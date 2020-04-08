package main

import "github.com/eagraf/synchronizer/workers"

// Messenger handles all communication
type Messenger struct {
	wm *workers.WorkerManager
}

// GetMessenger generates a new messenger singleton
// TODO ensure singularity
func GetMessenger(wm *workers.WorkerManager) *Messenger {
	m := Messenger{
		wm,
	}
	return &m
}

// Broadcast sends a message to each worker in separate parallel goroutines
func (m *Messenger) Broadcast() {

}
