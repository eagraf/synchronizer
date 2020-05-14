package messenger

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Test websocket client
type TestClient struct {
	conn *websocket.Conn
}

func newTestClient(URL string) (*TestClient, error) {
	var dialer websocket.Dialer
	parsedURL, _ := url.Parse(URL)
	fmt.Println("ws://" + parsedURL.Host)
	connection, _, err := dialer.Dial("ws://"+parsedURL.Host+"/", make(http.Header))
	if err != nil {
		return nil, err
	}

	tc := TestClient{
		conn: connection,
	}
	return &tc, nil
}

func (tc *TestClient) send(message *Message) {

}

func (tc *TestClient) receive() *Message {
	return nil
}

type mockService struct {
	connectionManager *ConnectionManager
	server            *httptest.Server
}

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
	ms.connectionManager.AddConnection(uuid.New().String(), r)

	// Temporarily promote request
	_, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
}

func TestMockServiceHandshake(t *testing.T) {
	ms := startMockService()
	_, err := newTestClient(ms.server.URL)
	if err != nil {
		t.Error(err.Error())
	}
}
