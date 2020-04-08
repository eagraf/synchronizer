package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestRegisterRemoveWorker(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(RegisterRoutes())
	port := ts.URL[17:]
	defer ts.Close()

	// Build the post request
	url := url.URL{Scheme: "ws", Host: "127.0.0.1:" + port, Path: "/workers", RawQuery: "workertype=cloudworker"}

	// Dial the websocket connection
	c, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		t.Errorf("Failed to dial %v", err)
	}
	defer c.Close()

	// Read the uuid returned after initial connection
	_, uuid, err := c.ReadMessage()
	if err != nil {
		t.Errorf("Failed to read %v", err)
	}

	if len(uuid) != 36 {
		t.Errorf("Invalid UUID returned")
	}

	// Build the delete request
	req, err := http.NewRequest("DELETE", ts.URL+"/workers/"+string(uuid)+"/", bytes.NewBuffer([]byte{}))
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
