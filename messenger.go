package main

// Messenger handles all communication
type Messenger struct {
	wm *WorkerManager
}

// GetMessenger generates a new messenger singleton
// TODO ensure singularity
func GetMessenger(wm *WorkerManager) *Messenger {
	m := Messenger{
		wm,
	}
	return &m
}

// Broadcast sends a message to each worker in separate parallel goroutines
func (m *Messenger) Broadcast() {

}
