package messenger

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"encoding/json"
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

func (tc *TestClient) send(message *Message) error {
	w, err := tc.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	zw := zlib.NewWriter(w)

	// Write offset
	offset := make([]byte, 4)
	binary.LittleEndian.PutUint32(offset, uint32(message.offset))
	zw.Write(offset)

	// Write metadata
	marshalled, err := json.Marshal(message.metadata)
	if err != nil {
		return err
	}

	zw.Write(marshalled)

	// Write payload
	if message.metadata.HasPayload {
		zw.Write(message.payload)
	}
	err = zw.Close()
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return nil
}

/*func (tc *TestClient) receive() (*Message, error) {

	_, buffer, err := tc.conn.ReadMessage()
	if err != nil {
		return err
	}
	message := FromBuffer(buffer)
	return message, nil
}*/

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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	_, msg, _ := conn.ReadMessage()
	fmt.Println(string(msg))

	zr, err := zlib.NewReader(bytes.NewReader(msg))
	if err != nil {
		fmt.Println("Failed to decompress: " + err.Error())
	}
	// Read into byte array
	inflated := new(bytes.Buffer)
	_, err = inflated.ReadFrom(zr)
	if err != nil {
		fmt.Println("Failed to decompress: " + err.Error())
	}
	fmt.Println(string(inflated.Bytes()))
}

func TestMockServiceHandshake(t *testing.T) {
	ms := startMockService()
	tc, err := newTestClient(ms.server.URL)
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
