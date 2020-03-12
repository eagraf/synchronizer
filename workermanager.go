package main

import (
	"errors"
	"net"

	"github.com/google/uuid"
)

// Worker represents a single worker device that executes tasks
// TODO make worker representation possible for heterogeneous hardware (i.e. phones + servers)
type Worker struct {
	UUID    string
	address net.IP
}

// WorkerManager keeps a table of all active wrokers
type WorkerManager struct {
	Workers map[string]Worker
}

// GetWorkerManager returns an instance of a workermanager
func GetWorkerManager() *WorkerManager {
	wm := WorkerManager{
		Workers: make(map[string]Worker),
	}
	return &wm
}

// AddWorker adds a worker to the workerset and returns the worker's uuid
func (wm *WorkerManager) AddWorker(address net.IP) string {
	var worker Worker
	// Generate new UUID for worker
	worker.UUID = uuid.New().String()
	worker.address = address

	// Add to workerset
	wm.Workers[worker.UUID] = worker

	return worker.UUID
}

// RemoveWorker removes a woker from the workerset
func (wm *WorkerManager) RemoveWorker(UUID string) error {
	_, ok := wm.Workers[UUID]
	if !ok {
		return errors.New("No worker with this UUID exists")
	}

	delete(wm.Workers, UUID)
	return nil
}
