package selector

import (
	"net/http"
	"testing"
)

func TestNewSelector(t *testing.T) {
	_, err := newSelector(3000, 3001)
	if err != nil {
		t.Error("No error expected from starting selector: " + err.Error())
	}

	// Test that external API started
	req, err := http.NewRequest("GET", "http://localhost:3000", nil)
	if err != nil {
		t.Error("Error constructing test API request: " + err.Error())
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Request failed: " + err.Error())
	}

	// Test that RPC started
	req, err = http.NewRequest("GET", "http://localhost:3001", nil)
	if err != nil {
		t.Error("Error constructing test API request: " + err.Error())
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Request failed: " + err.Error())
	}
}
