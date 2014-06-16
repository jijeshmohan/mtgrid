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
	for {
		messageType, p, err := ws.ReadMessage()
		ws.WriteMessage(messageType, p)
		if err != nil {
			log.Println("ERROR :", err)
			return
		}
	}
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
