package routes

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var ActiveClients = make(map[ClientConn]int)
var ActiveClientsRWMutex sync.RWMutex

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  net.Addr
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func addClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	ActiveClients[cc] = 0
	ActiveClientsRWMutex.Unlock()
}

func deleteClient(cc ClientConn) {
	ActiveClientsRWMutex.Lock()
	delete(ActiveClients, cc)
	ActiveClientsRWMutex.Unlock()
}

func broadcastMessage(messageType int, message []byte) {
	log.Println("Broadcasting")
	ActiveClientsRWMutex.RLock()
	defer ActiveClientsRWMutex.RUnlock()

	for client, _ := range ActiveClients {
		if err := client.websocket.WriteMessage(messageType, message); err != nil {
			return
		}
	}
}

func Socket(w http.ResponseWriter, r *http.Request) {
	log.Println("New connection")

	ws, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("Not a websocket handshake")
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	client := ws.RemoteAddr()
	sockCli := ClientConn{ws, client}
	addClient(sockCli)

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			deleteClient(sockCli)
			log.Println(err)
			return
		}
		broadcastMessage(messageType, p)
	}
}
