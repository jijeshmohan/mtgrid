package routes

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jijeshmohan/mtgrid/models"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgradeWs = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type NodeMsg struct {
	Status  string
	Message string
}

type NodeConnection struct {
	ws   *websocket.Conn
	send chan []byte
}

func (c *NodeConnection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		h.broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *NodeConnection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *NodeConnection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func Node(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeWs.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	nodeMsg := NodeMsg{}
	ws.ReadJSON(&nodeMsg)

	if err = processClient(nodeMsg); err != nil {
		log.Println(err)
		ws.Close()
		return
	}

	defer updateStatus(nodeMsg.Message, "Disconnected")

	c := &NodeConnection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump()
}

func updateStatus(devicename string, status string) error {
	device, err := models.GetDeviceWithName(strings.ToLower(devicename))
	if err != nil || device == nil {
		log.Println(err)
		return fmt.Errorf("Unable to find device : %s", devicename)
	}
	err = models.UpdateDeviceStatus(device, status)
	return err
}

func processClient(msg NodeMsg) error {
	if strings.ToLower(msg.Status) != "connected" {
		return fmt.Errorf("Invalid protocol from client")
	}

	err := updateStatus(msg.Message, "Connected")
	return err
}
