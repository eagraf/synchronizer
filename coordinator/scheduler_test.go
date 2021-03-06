package coordinator

import (
	"testing"

	"github.com/eagraf/synchronizer/service"
)

// Generic tests that all schedulers must pass

// Tests specific to schedulers
func TestNaiveSchedulerEqual(t *testing.T) {

	taskQueue := []*Task{
		{JobUUID: "1", TaskIndex: 1},
		{JobUUID: "2", TaskIndex: 1},
		{JobUUID: "3", TaskIndex: 1},
		{JobUUID: "4", TaskIndex: 1},
		{JobUUID: "5", TaskIndex: 1},
	}
	workerQueue := []*service.WorkersResponse_Worker{
		{WorkerUUID: "1"},
		{WorkerUUID: "2"},
		{WorkerUUID: "3"},
		{WorkerUUID: "4"},
		{WorkerUUID: "5"},
	}

	ns := naiveScheduler{}
	sched := ns.scheduleWorkers(taskQueue, workerQueue)

	for workerUUID, worker := range sched.assignments {
		if len(worker) != 1 {
			t.Error("Worker does not have the correct number of assigned jobs")
		}
		for jobUUID, job := range worker {
			if workerUUID != jobUUID {
				t.Error("Incorrect job assigned to worker")
			}
			if len(job) != 1 {
				t.Errorf("Worker has incorrect number of tasks from this job: %d", len(job))
			}
		}
	}
}

func TestNaiveSchedulerMoreTasks(t *testing.T) {
	taskQueue := []*Task{
		{JobUUID: "1", TaskIndex: 1},
		{JobUUID: "2", TaskIndex: 1},
		{JobUUID: "3", TaskIndex: 1},
		{JobUUID: "4", TaskIndex: 1},
		{JobUUID: "5", TaskIndex: 1},
	}
	workerQueue := []*service.WorkersResponse_Worker{
		{WorkerUUID: "1"},
		{WorkerUUID: "2"},
	}

	ns := naiveScheduler{}
	sched := ns.scheduleWorkers(taskQueue, workerQueue)

	if len(sched.assignments["1"]) != 3 {
		t.Error("Incorrect number of tasks assigned to first worker")
	}
	if len(sched.assignments["2"]) != 2 {
		t.Error("Incorrect number of tasks assigned to second worker")
	}
}

func TestNaiveSchedulerMoreWorkers(t *testing.T) {
	taskQueue := []*Task{
		{JobUUID: "1", TaskIndex: 1},
		{JobUUID: "2", TaskIndex: 1},
		{JobUUID: "3", TaskIndex: 1},
	}
	workerQueue := []*service.WorkersResponse_Worker{
		{WorkerUUID: "1"},
		{WorkerUUID: "2"},
		{WorkerUUID: "3"},
		{WorkerUUID: "4"},
		{WorkerUUID: "5"},
	}

	ns := naiveScheduler{}
	sched := ns.scheduleWorkers(taskQueue, workerQueue)

	if len(sched.assignments["1"]) != 1 {
		t.Error("Incorrect number of assignments for worker 1")
	}
	if len(sched.assignments["3"]) != 1 {
		t.Error("Incorrect number of assignments for worker 3")
	}
	if len(sched.unassignedWorkers) != 2 {
		t.Error("Incorrect number of unassigned workers")
	}
}

func TestNaiveScheduleDataServers(t *testing.T) {
	jobs := []*MapReduceJob{
		{JobType: "1"},
		{JobType: "2"},
		{JobType: "3"},
		{JobType: "4"},
		{JobType: "5"},
	}

	ds := []*dataServer{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	ns := naiveScheduler{}
	sched := ns.scheduleDataServers(jobs, ds)

	if len(sched.assignments) != 3 {
		t.Errorf("Incorrect number of assignments: %d", len(sched.assignments))
	}
	if len(sched.assignments["1"]) != 2 {
		t.Error("DataServer 1 has incorrect number of assignments")
	}
	if len(sched.assignments["3"]) != 1 {
		t.Error("DataServer 2 has incorrect number of assignments")
	}
}
func TestNaiveScheduleAggregators(t *testing.T) {
	jobs := []*MapReduceJob{
		{JobType: "1"},
		{JobType: "2"},
		{JobType: "3"},
		{JobType: "4"},
		{JobType: "5"},
	}

	ag := []*aggregator{
		{ID: "1"},
		{ID: "2"},
		{ID: "3"},
	}

	ns := naiveScheduler{}
	sched := ns.scheduleAggregators(jobs, ag)

	if len(sched.assignments) != 3 {
		t.Errorf("Incorrect number of assignments: %d", len(sched.assignments))
	}
	if len(sched.assignments["1"]) != 2 {
		t.Error("DataServer 1 has incorrect number of assignments")
	}
	if len(sched.assignments["3"]) != 1 {
		t.Error("DataServer 2 has incorrect number of assignments")
	}
}
