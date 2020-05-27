package messenger

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// TestClient represents the worker end of the websocket connection
// Test websocket client
type TestClient struct {
	conn *websocket.Conn
}

// NewTestClient creates a new test client
// Connect via websocket to mock service
func NewTestClient(URL string, clientID string) (*TestClient, error) {
	var dialer websocket.Dialer
	parsedURL, _ := url.Parse(URL)
	header := make(http.Header)
	header.Add("clientID", clientID)
	connection, _, err := dialer.Dial("ws://"+parsedURL.Host+parsedURL.EscapedPath(), header)
	if err != nil {
		return nil, err
	}

	tc := TestClient{
		conn: connection,
	}
	return &tc, nil
}

// Send message to mock service
func (tc *TestClient) Send(message *Message) error {
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
func (tc *TestClient) Receive() (*Message, error) {
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

// Close websocket
func (tc *TestClient) Close() error {
	//return tc.conn.Close()
	mb := MessageBuilder{}
	m, _ := mb.NewMessage(MessageClose, "").Done() // TODO ignore blank requestID

	w, err := tc.conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	// Send the message
	err = writeMessage(m, w)
	w.Close()
	if err != nil {
		return err
	}
	return nil

}
