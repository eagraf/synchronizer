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
	c := Coordinator{}
	writer := httptest.NewRecorder()
	c.createMapReduceJob(writer, request)

	// Check response out
	if writer.Code != 200 {
		t.Errorf("Incorrect code %d", writer.Code)
	}

	// Check coordinator state afterward

}
