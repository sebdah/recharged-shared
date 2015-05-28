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

	return
}

// Handler registering connections
func (this *Server) Handler(w http.ResponseWriter, r *http.Request) {
	var err error
	this.Conn, err = this.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("Error upgrading websocket connection: %s\n", err)
		return
	}

	// Register ping handler
	this.Conn.SetPingHandler(func(message string) error {
		log.Debug("Received ping message")
		this.Conn.WriteMessage(websocket.PongMessage, []byte(message))
		return nil
	})

	// Instanciate a new communicator
	communicator := NewCommunicator(this.Conn)
	log.Debug("Starting websockets communication channel")
	go communicator.Reader(this.ReadMessage)
	go communicator.Writer(this.WriteMessage)
}
