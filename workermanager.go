package main

import (
	"errors"
	"net"
	"sync"

	"github.com/google/uuid"
)

// Worker represents a single worker device that executes tasks
// TODO make worker representation possible for heterogeneous hardware (i.e. phones + servers)
// TODO enforce enumerated deviceTypes (i.e. "cloud", "mobile" workers)
type Worker struct {
	UUID       string
	address    net.IP
	workerType string
}

// WorkerManager keeps a table of all active wrokers
type WorkerManager struct {
	Workers     map[string]Worker
	workerMutex sync.RWMutex
}

// WorkerManager singleton
var wmSingleton = WorkerManager{
	Workers:     make(map[string]Worker),
	workerMutex: sync.RWMutex{},
}

// GetWorkerManager returns an instance of a workermanager
// TODO handle singularity
func GetWorkerManager() *WorkerManager {
	return &wmSingleton
}

// AddWorker adds a worker to the workerset and returns the worker's uuid
func (wm *WorkerManager) AddWorker(address net.IP, workerType string) string {
	var worker Worker
	// Generate new UUID for worker
	worker.UUID = uuid.New().String()
	worker.address = address
	worker.workerType = workerType

	// Get write mutex
	wm.workerMutex.Lock()
	// Add to workerset
	wm.Workers[worker.UUID] = worker
	// Release write mutex
	wm.workerMutex.Unlock()

	return worker.UUID
}

// RemoveWorker removes a woker from the workerset
func (wm *WorkerManager) RemoveWorker(UUID string) error {
	// Get read mutex
	wm.workerMutex.RLock()
	// Check if UUID exists
	_, ok := wm.Workers[UUID]
	if !ok {
		return errors.New("No worker with this UUID exists")
	}
	// Release read mutex
	wm.workerMutex.RUnlock()

	// Delete the worker
	delete(wm.Workers, UUID)
	return nil
}
