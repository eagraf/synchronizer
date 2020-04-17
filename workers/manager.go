package workers

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/eagraf/synchronizer/tasks"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Worker represents a single worker device that executes tasks
// TODO make worker representation possible for heterogeneous hardware (i.e. phones + servers)
// TODO enforce enumerated deviceTypes (i.e. "cloud", "mobile" workers)
type Worker struct {
	UUID       string
	workerType string
	connection *websocket.Conn
}

// WorkerManager keeps a table of all active wrokers
type WorkerManager struct {
	Workers          map[string]Worker   // Table of all workers that are available or working
	AvailableWorkers chan Worker         // Workers that are looking for tasks
	AllocatedWorkers map[string][]Worker // Table of workers assigned to specific tasks, key is a taskUUID
	MapTaskQueue     chan *tasks.Intent  // Cue of incoming map task intents
	workersMutex     sync.RWMutex
	allocationMutex  sync.RWMutex
}

// WorkerManager singleton
var wmSingleton = WorkerManager{
	Workers:          make(map[string]Worker),
	AvailableWorkers: make(chan Worker, 1024),
	AllocatedWorkers: make(map[string][]Worker),
	MapTaskQueue:     make(chan *tasks.Intent),
	workersMutex:     sync.RWMutex{},
	allocationMutex:  sync.RWMutex{},
}

// GetWorkerManager returns an instance of a workermanager
// TODO handle singularity
func GetWorkerManager() *WorkerManager {
	return &wmSingleton
}

// AddWorker adds a worker to the workerset and returns the worker's uuid
func (wm *WorkerManager) AddWorker(workerType string, connection *websocket.Conn) string {
	var worker Worker
	// Generate new UUID for worker
	worker.UUID = uuid.New().String()
	worker.workerType = workerType
	worker.connection = connection

	// Get write mutex
	wm.workersMutex.Lock()
	// Add to workerset
	wm.Workers[worker.UUID] = worker
	// Add to available workers
	wm.AvailableWorkers <- worker
	// Release write mutex
	wm.workersMutex.Unlock()

	return worker.UUID
}

// RemoveWorker removes a woker from the workerset
func (wm *WorkerManager) RemoveWorker(UUID string) error {
	// Get read mutex
	wm.workersMutex.RLock()
	// Check if UUID exists
	_, ok := wm.Workers[UUID]
	if !ok {
		return errors.New("No worker with this UUID exists")
	}
	// Release read mutex
	wm.workersMutex.RUnlock()

	// Delete the worker
	delete(wm.Workers, UUID)
	return nil
}

// Start listening for incoming map tasks originated by the scheduler
func (wm *WorkerManager) Start() {
	// Listen for new intents and workers
	// TODO verify this setup doesn't lead to excesive waiting / explore parallelization
	//
	// Currently just a naive implementation that takes the first available worker
	// Future extensions:
	//  - Account for device capabilities (bandwidth, compute, GPU, etc.)
	//  - Worker management is mission critical, and should be replicated using Paxos or equivalent
	//  - Support multischeduling workers
	go func() {
		for {
			select {
			case mapIntent := <-wm.MapTaskQueue:
				fmt.Println("Received map intent")
				// Get write mutex
				wm.allocationMutex.Lock()
				worker := <-wm.AvailableWorkers
				// TODO need to check if worker is still available
				fmt.Printf("Allocating worker %v to task %v (intent listener)\n", worker.UUID, mapIntent.TaskUUID)
				// Release write mutex
				wm.allocationMutex.Unlock()
				go wm.MessageWorker(worker.UUID, mapIntent.TaskUUID, mapIntent)

			case worker := <-wm.AvailableWorkers:
				fmt.Println("Received worker")
				// Get write mutex
				wm.allocationMutex.Lock()
				mapIntent := <-wm.MapTaskQueue
				// TODO need to check if worker is still available
				fmt.Printf("Allocating worker %v to task %v (worker listener)\n", worker.UUID, mapIntent.TaskUUID)
				// Release write mutex
				wm.allocationMutex.Unlock()
				go wm.MessageWorker(worker.UUID, mapIntent.TaskUUID, mapIntent)
			}
		}
	}()
}

// MessageWorker sends a map intent to a worker
// TODO determine if WriteCloser can be fixed to work with gobs => could prevent one buffer copy
func (wm *WorkerManager) MessageWorker(workerUUID, taskUUID string, mapIntent *tasks.Intent) error {
	// Decompose the map intent into a JSON like map
	buffer, err := json.Marshal(mapIntent)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	// Send the intent to the worker
	err = wm.Workers[workerUUID].connection.WriteMessage(websocket.TextMessage, buffer)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
