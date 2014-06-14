package routes

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jijeshmohan/mtgrid/models"
)

type NodeMsg struct {
	Status  string
	Message string
}

var upgradeWs = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Node(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeWs.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("Not a websocket handshake")
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
	for {
		messageType, p, err := ws.ReadMessage()
		ws.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func processClient(msg NodeMsg) error {
	if strings.ToLower(msg.Status) != "connected" {
		return fmt.Errorf("Invalid protocol from client")
	}
	device, err := models.GetDeviceWithName(strings.ToLower(msg.Message))
	if err != nil || device == nil {
		log.Println(err)
		return fmt.Errorf("Unable to find device : %s", msg.Message)
	}
	err = models.UpdateDeviceStatus(device, "Connected")
	return err
}
