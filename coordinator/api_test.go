package coordinator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func makeCreateMapReduceRequest(c *Coordinator) (*httptest.ResponseRecorder, error) {
	// Mock request
	reqBody, err := json.Marshal(MapReduceJob{
		JobType:    "TestTask",
		TaskSize:   1000,
		TaskNumber: 1000,
	})
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", "http://localhost:2216/jobs/", bytes.NewBuffer(reqBody))
	request.Header.Set("Content-type", "application/json")
	if err != nil {
		return nil, err
	}

	// Pass mock request into api call
	writer := httptest.NewRecorder()
	c.createMapReduceJob(writer, request)

	return writer, nil
}

func TestCreateMapReduceRequest(t *testing.T) {
	/*
		// Mock request
		reqBody, err := json.Marshal(MapReduceJob{
			JobType:    "TestTask",
			TaskSize:   1000,
			TaskNumber: 1000,
		})
		if err != nil {
			t.Error(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:2216/jobs/", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-type", "application/json")
		if err != nil {
			t.Error(err)
		}

		// Pass mock request into api call
		c := Coordinator{
			activeJobs: make(map[string]*MapReduceJob),
			taskQueue:  make([]*Task, 0, 1<<20), // Give a rather large initial capacity
		}
		writer := httptest.NewRecorder()
		c.createMapReduceJob(writer, request)*/

	c := &Coordinator{
		activeJobs: make(map[string]*MapReduceJob),
		taskQueue:  make([]*Task, 0, 1<<20), // Give a rather large initial capacity
	}
	writer, err := makeCreateMapReduceRequest(c)
	if err != nil {
		t.Error(err.Error())
	}

	// Check response out
	if writer.Code != 200 {
		t.Errorf("Incorrect code %d", writer.Code)
	}

	// Check coordinator state afterward
	if len(c.taskQueue) != 1000 {
		t.Errorf("Coordinator task queue is not len 1000: %d ", len(c.taskQueue))
	}
	if len(c.activeJobs) != 1 {
		t.Errorf("Coordinator does not have the right number of active jobs %d", len(c.activeJobs))
	}

	// Check response body
	buf := writer.Body.Bytes()
	var res MapReduceJobResponse
	err = json.Unmarshal(buf, &res)
	if err != nil {
		t.Error(err.Error())
	}

	if len(res.JobUUID) != 36 {
		t.Error("Invalid UUID")
	}
}

// Run this test with go test -race -run TestTaskQueueWait
func TestTaskQueueRace(t *testing.T) {
	c := &Coordinator{
		activeJobs: make(map[string]*MapReduceJob),
		taskQueue:  make([]*Task, 0, 1<<20), // Give a rather large initial capacity
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		makeCreateMapReduceRequest(c)
	}()
	go func() {
		defer wg.Done()
		makeCreateMapReduceRequest(c)
	}()
	wg.Wait()
}
