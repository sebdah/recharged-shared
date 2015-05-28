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
	PingMessage     chan string
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
	client.PingMessage = make(chan string)

	client.connect()

	return
}

// Connect to a websockets server
func (this *Client) connect() {
	rawConn, err := net.Dial("tcp", this.Endpoint.Host)
	if err != nil {
		panic(err)
	}

	this.Conn, _, err = websocket.NewClient(
		rawConn,
		this.Endpoint,
		this.Headers,
		this.ReadBufferSize,
		this.WriteBufferSize)
	if err != nil {
		panic(err)
	}

	// Set some limits
	this.Conn.SetReadLimit(maxMessageSize)
	this.Conn.SetReadDeadline(time.Now().Add(pongWait))

	log.Info("Connected to endpoint '%s' via websockets\n", this.Endpoint)

	// Instanciate a new communicator
	communicator := NewCommunicator(this.Conn)
	log.Debug("Starting websockets communication channel")
	go communicator.Reader(this.ReadMessage)
	go communicator.Writer(this.WriteMessage, this.PingMessage)
}

// Set ping handler
func (this *Client) SetPingHandler(h func(string) error) {
	log.Debug("Setting new ping handler")
	this.Conn.SetPingHandler(h)
}

// Set pong handler
func (this *Client) SetPongHandler(h func(string) error) {
	log.Debug("Setting new pong handler")
	this.Conn.SetPongHandler(h)
}
