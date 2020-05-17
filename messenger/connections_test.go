package messenger

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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
	connectionManager *connectionManager
	server            *httptest.Server
	t                 *testing.T
	onReceives        int
	onSends           int
	onCloses          int
}

// Start the test server
func startMockService(t *testing.T) *mockService {
	ps := newPubSub()
	cm := newConnectionManager(ps)
	ms := mockService{
		connectionManager: cm,
		t:                 t,
		onReceives:        0,
		onSends:           0,
		onCloses:          0,
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
	err := ms.connectionManager.AddConnection(r.Header.Get("clientID"), w, r)
	if err != nil {
		ms.t.Error(err.Error())
	}
	err = ms.connectionManager.subscriptions.AddSubscription(r.Header.Get("clientID"), ms)
	if err != nil {
		ms.t.Error(err.Error())
	}
}

func (ms *mockService) GetIdentifier() string {
	return "mock-service-identifier"
}

func (ms *mockService) OnReceive(topic string, message *Message) {
	ms.onReceives++
}

func (ms *mockService) OnSend(topic string, message *Message) {
	ms.onSends++
}

func (ms *mockService) OnClose(topic string) {
	ms.onCloses++
}

// Tests AddConnection
func TestMockServiceHandshake(t *testing.T) {
	ms := startMockService(t)
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

	// TODO figure out a better way of testing that doesn't require sleeps to enforce sequential=>deterministic testing
	time.Sleep(time.Second)

	ms.connectionManager.RemoveConnection("client-1")
	if len(ms.connectionManager.connections) != 0 {
		t.Error("There should be no remaining connections")
	}

	// These tests must be performed at end to avoid race condition
	if ms.onReceives != 1 {
		fmt.Println(ms.onReceives)
		t.Error("Expected 1 invocation of OnReceive")
	}
	if ms.onCloses != 1 {
		t.Error("Expected 1 invocation of OnClose")
	}
}

func TestMessaging(t *testing.T) {
	ms := startMockService(t)
	tc, err := newTestClient(ms.server.URL, "client-1")
	if err != nil {
		t.Error(err.Error())
	}

	if len(ms.connectionManager.connections) != 1 {
		t.Error("There should be one connection")
	}

	// Construct message
	mb := MessageBuilder{}
	m, _ := mb.NewMessage("test-message", "request-id").Done()

	// Test sending from sychronizer to the testClient
	ms.connectionManager.Send("client-1", m)
	_, err = tc.receive()
	if err != nil {
		t.Error(err.Error())
	}

	// testClient sends response
	m2, _ := mb.NewMessage("response-message", "request-id").Done()
	err = tc.send(m2)
	if err != nil {
		t.Error(err.Error())
	}
	time.Sleep(time.Second)

	ms.connectionManager.RemoveConnection("client-1")
	if len(ms.connectionManager.connections) != 0 {
		t.Error("There should be no remaining connections")
	}

	// These tests must be performed at end to avoid race condition
	if ms.onReceives != 1 {
		t.Error("Expected 1 invocation of OnReceive")
	}
	if ms.onSends != 1 {
		t.Error("Expected 1 invocation of OnSends")
	}
	if ms.onCloses != 1 {
		t.Error("Expected 1 invocation of OnClose")
	}
}

func TestConcurrentSend(t *testing.T) {
	ms := startMockService(t)
	tc, err := newTestClient(ms.server.URL, "client-1")
	if err != nil {
		t.Error(err.Error())
	}

	if len(ms.connectionManager.connections) != 1 {
		t.Error("There should be one connection")
	}

	// Construct message
	mb := MessageBuilder{}
	m, _ := mb.NewMessage("test-message", "request-id").Done()

	// Test sending from sychronizer to the testClient
	ms.connectionManager.Send("client-1", m)
	ms.connectionManager.Send("client-1", m)
	_, err = tc.receive()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = tc.receive()
	if err != nil {
		t.Error(err.Error())
	}

	// testClient sends response
	m2, _ := mb.NewMessage("response-message", "request-id").Done()
	err = tc.send(m2)
	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(time.Second)

	ms.connectionManager.RemoveConnection("client-1")
	if len(ms.connectionManager.connections) != 0 {
		t.Error("There should be no remaining connections")
	}
	// These tests must be performed at end to avoid race condition
	if ms.onReceives != 1 {
		t.Error("Expected 1 invocation of OnReceive")
	}
	if ms.onSends != 2 {
		t.Error("Expected 1 invocation of OnSends")
	}
	if ms.onCloses != 1 {
		t.Error("Expected 1 invocation of OnClose")
	}
}

func TestMultipleConnections(t *testing.T) {
	ms := startMockService(t)
	tc1, err := newTestClient(ms.server.URL, "client-1")
	if err != nil {
		t.Error(err.Error())
	}
	tc2, err := newTestClient(ms.server.URL, "client-2")
	if err != nil {
		t.Error(err.Error())
	}

	if len(ms.connectionManager.connections) != 2 {
		t.Error("There should be two connections")
	}

	// Construct message
	mb := MessageBuilder{}
	m, _ := mb.NewMessage("test-message", "request-id").Done()

	ms.connectionManager.Send("client-1", m)
	_, err = tc1.receive()
	if err != nil {
		t.Error(err.Error())
	}
	ms.connectionManager.Send("client-2", m)
	_, err = tc2.receive()
	if err != nil {
		t.Error(err.Error())
	}

	time.Sleep(time.Second)

	ms.connectionManager.RemoveConnection("client-1")
	if len(ms.connectionManager.connections) != 1 {
		t.Error("There should be only one connection")
	}
	ms.connectionManager.RemoveConnection("client-2")
	if len(ms.connectionManager.connections) != 0 {
		t.Error("There should be no connections")
	}
	// These tests must be performed at end to avoid race condition
	if ms.onReceives != 0 {
		t.Error("Expected 1 invocation of OnReceive")
	}
	if ms.onSends != 2 {
		t.Error("Expected 1 invocation of OnSends")
	}
	if ms.onCloses != 2 {
		t.Error("Expected 1 invocation of OnClose")
	}
}
