package coordinator

import "github.com/eagraf/synchronizer/service"

type naiveScheduler struct{}

func (ns *naiveScheduler) scheduleWorkers(taskQueue []*Task, workerQueue []*service.WorkersResponse_Worker) *workerSchedule {
	res := &workerSchedule{
		assignments:       make(map[string]map[string][]int),
		unassignedWorkers: make(map[string]bool),
	}
	for i, task := range taskQueue {
		worker := workerQueue[i%len(workerQueue)]

		if _, ok := res.assignments[worker.WorkerUUID]; ok == false {
			res.assignments[worker.WorkerUUID] = make(map[string][]int)
		}
		if _, ok := res.assignments[worker.WorkerUUID][task.JobUUID]; ok == false {
			res.assignments[worker.WorkerUUID][task.JobUUID] = make([]int, 0)
		}

		res.assignments[worker.WorkerUUID][task.JobUUID] = append(res.assignments[worker.WorkerUUID][task.JobUUID], task.TaskIndex)
	}

	// Populate unassigned workers
	for i := len(taskQueue); i < len(workerQueue); i++ {
		res.unassignedWorkers[workerQueue[i].WorkerUUID] = true
	}
	return res
}

func (ns *naiveScheduler) scheduleDataServers(jobs []*MapReduceJob, dataServers []*dataServer) *dataServerSchedule {
	res := &dataServerSchedule{
		assignments: make(map[string][]string),
	}

	// Initialize each data server in assignments map
	for _, ds := range dataServers {
		res.assignments[ds.ID] = make([]string, 0)
	}
	// Populate assignments
	for i, job := range jobs {
		dsUUID := dataServers[i%len(dataServers)].ID
		res.assignments[dsUUID] = append(res.assignments[dsUUID], job.JobUUID)
	}
	return res
}

func (ns *naiveScheduler) scheduleAggregators(jobs []*MapReduceJob, aggregators []*aggregator) *aggregatorSchedule {
	res := &aggregatorSchedule{
		assignments: make(map[string][]string),
	}

	// Initialize each data server in assignments map
	for _, ds := range aggregators {
		res.assignments[ds.ID] = make([]string, 0)
	}
	// Populate assignments
	for i, job := range jobs {
		agUUID := aggregators[i%len(aggregators)].ID
		res.assignments[agUUID] = append(res.assignments[agUUID], job.JobUUID)
	}
	return res
}
