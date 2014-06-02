package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Node(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	// client := ws.RemoteAddr()

	for {
		messageType, p, err := ws.ReadMessage()
		ws.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
