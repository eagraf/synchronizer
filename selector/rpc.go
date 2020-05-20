package selector

// GetWorkers RPC interface
func (s *Selector) GetWorkers(request WorkerRequest, reply *WorkerResponse) error {
	// Convert to slice
	res := make([]Worker, 0, len(s.workers))
	for _, w := range s.workers {
		res = append(res, *w)
	}
	reply.workers = res
	return nil
}

func (s *Selector) Handoff(request HandoffRequest, reply *HandoffRequest) error {
	return nil
}
