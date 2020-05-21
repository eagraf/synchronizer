package selector

// RPCHandler holds all exported RPC methods
type RPCHandler struct {
	selector *Selector
}

// WorkerRequest is input for GetWorkers RPC call
type WorkerRequest struct {
	numRequested int
}

// WorkerResponse is output for GetWorkers RPC call
type WorkerResponse struct {
	workers []Worker
}

// GetWorkers RPC interface
func (rh *RPCHandler) GetWorkers(request WorkerRequest, reply *WorkerResponse) error {
	// Convert to slice
	res := make([]Worker, 0, len(rh.selector.workers))
	for _, w := range rh.selector.workers {
		res = append(res, *w)
	}
	reply.workers = res
	return nil
}

// HandoffRequest is input for Handoff RPC call
type HandoffRequest struct {
}

// HandoffResponse is output for Handoff RPC call
type HandoffResponse struct {
}

// Handoff RPC interface
func (rh *RPCHandler) Handoff(request HandoffRequest, reply *HandoffRequest) error {
	return nil
}
