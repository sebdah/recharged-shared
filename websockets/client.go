package websockets

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn            *websocket.Conn
	Endpoint        *url.URL
	Headers         http.Header
	ReadBufferSize  int
	WriteBufferSize int
	WriteMessage    chan string
	ReadMessage     chan string
	PingHandler     func(string) error
	PongHandler     func(string) error
}

// Constructor
func NewClient(endpoint *url.URL) (client *Client) {
	client = new(Client)
	client.Endpoint = endpoint
	client.ReadBufferSize = 1024
	client.WriteBufferSize = 1024
	client.Headers = http.Header{
		"Origin":                   {fmt.Sprintf("%s://%s", endpoint.Scheme, endpoint.Host)},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
	}

	// Create channels
	client.ReadMessage = make(chan string)
	client.WriteMessage = make(chan string)

	// Setup default ping and pong handlers
	client.PingHandler = func(m string) (err error) {
		log.Debug("Received ping: %s", m)
		client.Conn.WriteMessage(websocket.PongMessage, []byte{})
		return
	}
	client.PongHandler = func(m string) (err error) {
		log.Debug("Received pong: %s", m)
		return
	}

	client.connect()

	return
}

// Connect to a websockets server
func (this *Client) connect() {
	rawConn, err := net.Dial("tcp", this.Endpoint.Host)
	if err != nil {
		panic(err)
	}

	conn, _, err := websocket.NewClient(
		rawConn,
		this.Endpoint,
		this.Headers,
		this.ReadBufferSize,
		this.WriteBufferSize)
	if err != nil {
		panic(err)
	}
	this.Conn = conn

	// Set some limits
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))

	// Set the ping and pong handlers
	conn.SetPingHandler(this.PingHandler)
	conn.SetPongHandler(this.PongHandler)

	log.Info("Connected to endpoint '%s' via websockets\n", this.Endpoint)

	// Instanciate a new communicator
	communicator := NewCommunicator(conn)
	log.Debug("Starting websockets communication channel")
	go communicator.Reader(this.ReadMessage)
	go communicator.Writer(this.WriteMessage)
}
