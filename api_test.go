package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostWorker(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(RegisterRoutes())
	defer ts.Close()

	// Build the post request
	reqBody, _ := json.Marshal(map[string]interface{}{"ip": "192.0.2.1", "workerType": "cloud_worker"})
	req, err := http.NewRequest("POST", ts.URL+"/workers/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		t.Errorf("Bad POST response code %s", res.Status)
	}
}
