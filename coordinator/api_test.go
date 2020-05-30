package coordinator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateMapReduceRequest(t *testing.T) {

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
	c.createMapReduceJob(writer, request)

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
