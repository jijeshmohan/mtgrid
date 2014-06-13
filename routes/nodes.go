package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	// "github.com/jijeshmohan/mtgrid/models"
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
	// client := ws.RemoteAddr()
	nodeMsg := NodeMsg{}
	ws.ReadJSON(&nodeMsg)
	log.Println(nodeMsg)
	// TODO : process client connection
	for {
		messageType, p, err := ws.ReadMessage()
		ws.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
