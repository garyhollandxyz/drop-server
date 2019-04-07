package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Nexus struct {
	clients   map[*Client]bool
	broadcast chan ChatMessage
}

func newNexus() Nexus {
	return Nexus{
		clients:   make(map[*Client]bool),
		broadcast: make(chan ChatMessage),
	}
}

func (n Nexus) run() {
	http.HandleFunc("/ws", n.handleConnection)
}

func (n Nexus) handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("Made a connection")
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()
	client := &Client{conn: ws, nexus: n, send: make(chan ChatMessage)}
	n.clients[client] = true

	go client.streamIn()
	go client.streamOut()

	for {
		msg := <-n.broadcast

		for client := range n.clients {
			client.send <- msg
		}
	}
}
