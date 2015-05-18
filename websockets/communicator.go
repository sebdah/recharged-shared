package websockets

import (
	"github.com/gorilla/websocket"
)

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
		_, msg, err := this.conn.ReadMessage()
		if err != nil {
			return
		}

		// Send the received message to the channel
		c_recv <- string(msg)
	}
}

// Writer
func (this *Communicator) Writer(c_send chan string) {
	log.Debug("Write communicator started")
	send_msg := ""

	for {
		// Read the c_send channel
		send_msg = <-c_send
		if send_msg != "" {
			this.conn.WriteMessage(websocket.TextMessage, []byte(send_msg))
		}
	}
}
