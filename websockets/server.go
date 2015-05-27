package websockets

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	Conn            *websocket.Conn
	Upgrader        websocket.Upgrader
	ReadBufferSize  int
	WriteBufferSize int
	WriteMessage    chan string
	ReadMessage     chan string
	PingHandler     func(string) error
	PongHandler     func(string) error
}

func NewServer() (server *Server) {
	server = new(Server)
	server.ReadBufferSize = 1024
	server.WriteBufferSize = 1024

	server.Upgrader = websocket.Upgrader{
		ReadBufferSize:  server.ReadBufferSize,
		WriteBufferSize: server.WriteBufferSize,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	// Create channels
	server.ReadMessage = make(chan string)
	server.WriteMessage = make(chan string)

	// Setup default ping and pong handlers
	server.PingHandler = func(m string) (err error) {
		log.Debug("Received ping: %s", m)
		server.Conn.WriteMessage(websocket.PongMessage, []byte{})
		return
	}
	server.PongHandler = func(m string) (err error) {
		log.Debug("Received pong: %s", m)
		return
	}

	return
}

// Handler registering connections
func (this *Server) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := this.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Error upgrading websocket connection: %s\n", err)
		return
	}
	this.Conn = conn

	// Set the ping and pong handlers
	conn.SetPingHandler(this.PingHandler)
	conn.SetPongHandler(this.PongHandler)

	// Instanciate a new communicator
	communicator := NewCommunicator(conn)
	log.Debug("Starting websockets communication channel")
	go communicator.Reader(this.ReadMessage)
	go communicator.Writer(this.WriteMessage)
}
