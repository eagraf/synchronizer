package messenger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/websocket"
)

// TestClient represents the worker end of the websocket connection
// Test websocket client
type TestClient struct {
	conn *websocket.Conn
}

// Connect via websocket to mock service
func newTestClient(URL string, clientID string) (*TestClient, error) {
	var dialer websocket.Dialer
	parsedURL, _ := url.Parse(URL)
	header := make(http.Header)
	header.Add("clientID", clientID)
	connection, _, err := dialer.Dial("ws://"+parsedURL.Host+"/", header)
	if err != nil {
		return nil, err
	}

	tc := TestClient{
		conn: connection,
	}
	return &tc, nil
}

// Send message to mock service
func (tc *TestClient) send(message *Message) error {
	w, err := tc.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	// Send the message
	err = writeMessage(message, w)
	w.Close()
	if err != nil {
		return err
	}
	return nil
}

// Receive a message from the mock service
func (tc *TestClient) receive() (*Message, error) {
	_, buffer, err := tc.conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	message, err := readMessage(buffer)
	if err != nil {
		return nil, err
	}
	return message, nil
}

// mockService provides an HTTP endpoint for initiating test websocket connections
type mockService struct {
	connectionManager *ConnectionManager
	server            *httptest.Server
}

// Start the test server
func startMockService() *mockService {
	ps := newPubSub()
	cm := newConnectionManager(ps)
	ms := mockService{
		connectionManager: cm,
	}

	ms.server = httptest.NewServer(http.HandlerFunc(ms.mockWebsocketEndpoint))
	return &ms
}

// upgrader promotes a standard HTTP/HTTPS connection to a websocket connection
// TODO implement CheckOrigin
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(*http.Request) bool { return true },
}

// Mock websocket endpoint
func (ms *mockService) mockWebsocketEndpoint(w http.ResponseWriter, r *http.Request) {
	ms.connectionManager.AddConnection(r.Header.Get("clientID"), r)

	// Temporarily promote request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	_, buffer, _ := conn.ReadMessage()
	fmt.Println(string(buffer))

	msg, err := readMessage(buffer)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(msg.offset)
	fmt.Println(msg.metadata.Request)
}

// Tests AddConnection
func TestMockServiceHandshake(t *testing.T) {
	ms := startMockService()
	tc, err := newTestClient(ms.server.URL, "client-1")
	if err != nil {
		t.Error(err.Error())
	}

	// Test sending
	mb := MessageBuilder{}
	m, _ := mb.NewMessage("test-message", "test-request").Done()
	err = tc.send(m)
	if err != nil {
		t.Error(err.Error())
	}
}
