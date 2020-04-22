package workers

import (
	"errors"
	"fmt"
	"sync"

	"github.com/eagraf/synchronizer/messenger"

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
	connection *websocket.Conn `json:"-"`
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

// RegistrationResponse contains important metadata for a worker
type RegistrationResponse struct {
	UUID string
}

// WorkerManager singleton
var wmSingleton = WorkerManager{
	Workers:          make(map[string]Worker),
	AvailableWorkers: make(chan Worker, 1024),
	AllocatedWorkers: make(map[string][]Worker),      // TODO actually use allocations rather than relying on availability queue
	MapTaskQueue:     make(chan *tasks.Intent, 1024), // TODO figure out optimal buffering length
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
	// TODO determine if WLock / WUnlock is more proper
	wm.workersMutex.Lock()
	// Add to workerset
	wm.Workers[worker.UUID] = worker
	// Add to available workers
	wm.AvailableWorkers <- worker
	// Release write mutex
	wm.workersMutex.Unlock()

	// Register this worker with the messenger and send a registration response
	m := messenger.GetMessengerSingleton()
	m.AddConnection(worker.UUID, worker.connection)
	m.AddSubscriber(wm, []string{worker.UUID}) // Subscribe to worker
	m.SendMessage(worker.UUID, &RegistrationResponse{worker.UUID})

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
				wm.allocate(worker, mapIntent)
				wm.allocationMutex.Unlock()
				// Check if worker is still connected
				/*if _, ok := wm.Workers[worker.UUID]; ok == true {
					fmt.Printf("Allocating worker %v to task %v (intent listener)\n", worker.UUID, mapIntent.TaskUUID)
					// Release write mutex
					// Send message to worker and add subscription
					m := messenger.GetMessengerSingleton()
					ti, _ := tasks.GetTaskServiceSingleton().GetTaskInstance(mapIntent.TaskUUID)
					// TODO error handling
					// TODO better handling of subscribers
					m.AddSubscriber(ti, []string{worker.UUID})
					m.AddSubscriber(wm, []string{worker.UUID})
					m.SendMessage(worker.UUID, mapIntent)
				}*/

			case worker := <-wm.AvailableWorkers:
				fmt.Println("Received worker")
				// Get write mutex
				// TODO This mutex might need to surroudn ^^ worker channel, not mapTask channel
				wm.allocationMutex.Lock()
				mapIntent := <-wm.MapTaskQueue
				wm.allocate(worker, mapIntent)
				wm.allocationMutex.Unlock()
				/*
					// TODO need to check if worker is still available
					fmt.Printf("Allocating worker %v to task %v (worker listener)\n", worker.UUID, mapIntent.TaskUUID)
					// Release write mutex
					wm.allocationMutex.Unlock()

					// Send message to worker and add subscription
					m := messenger.GetMessengerSingleton()
					ti, _ := tasks.GetTaskServiceSingleton().GetTaskInstance(mapIntent.TaskUUID)
					// TODO error handling
					m.AddSubscriber(ti, []string{worker.UUID})
					m.AddSubscriber(wm, []string{worker.UUID})
					m.SendMessage(worker.UUID, mapIntent)*/
			}
		}
	}()
}

func (wm *WorkerManager) allocate(worker Worker, mapIntent *tasks.Intent) {
	// Check if worker is still connected
	if _, ok := wm.Workers[worker.UUID]; ok == true {
		fmt.Printf("Allocating worker %v to task %v (intent listener)\n", worker.UUID, mapIntent.TaskUUID)
		// Send message to worker and add subscription
		m := messenger.GetMessengerSingleton()
		ti, _ := tasks.GetTaskServiceSingleton().GetTaskInstance(mapIntent.TaskUUID)
		// TODO error handling
		// TODO better handling of subscribers
		m.AddSubscriber(ti, []string{worker.UUID})
		m.SendMessage(worker.UUID, mapIntent)
	}
}

// AddSubscription implements a messenger subscriber method
func (wm *WorkerManager) AddSubscription(topic string) {
	// Do nothing for now
}

// GetUUID implements a messenger subscriber method
func (wm *WorkerManager) GetUUID() string {
	return "WorkerManager" // Probably should define this as a constant somewhere
}

// OnReceive implements a messenger subscriber method
func (wm *WorkerManager) OnReceive(topic string, m *map[string]interface{}) {
	// TODO rn assuming topic is worker UUID
	// TODO evaluate if mutex is necessary here (in fact it often breaks things if you add it back)
	worker := wm.Workers[topic]
	//wm.allocationMutex.Lock()
	wm.AvailableWorkers <- worker
	//wm.allocationMutex.Unlock()
}

// OnClose implements a messenger subscriber method
func (wm *WorkerManager) OnClose(topic string) {
	// TODO this is Assuming topic is workerUUID
	// Remove worker
	fmt.Println("goodbye")
	wm.workersMutex.Lock()
	// Delete from workers
	fmt.Println("hello")
	if _, ok := wm.Workers[topic]; ok == true {
		fmt.Println("Hello" + topic)
		delete(wm.Workers, topic)
	}
	wm.workersMutex.Unlock()
}
