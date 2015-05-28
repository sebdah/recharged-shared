package websockets

import "github.com/gorilla/websocket"

type Communicator struct {
	conn *websocket.Conn
}

// Constructor
func NewCommunicator(conn *websocket.Conn) (communicator *Communicator) {
	communicator = new(Communicator)
	communicator.conn = conn

	return
}

// Reader function
func (this *Communicator) Reader(c_recv chan string) {
	log.Debug("Read communicator started")

	for {
		_, message, err := this.conn.ReadMessage()
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Send the received message to the channel
		c_recv <- string(message)
	}
}

// Writer
func (this *Communicator) Writer(c_send chan string) {
	log.Debug("Write communicator started")

	for {
		message, ok := <-c_send
		if !ok {
			log.Debug("Sending close message")
			this.conn.WriteMessage(websocket.CloseMessage, []byte{})
		} else {
			log.Debug("Sending text message")
			this.conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}

// Pinger
func (this *Communicator) Pinger(c_ping chan string) {
	log.Debug("Ping communicator started")

	for {
		message := <-c_ping
		log.Debug("Sending ping message")
		this.conn.WriteMessage(websocket.PingMessage, []byte(message))
	}
}
