// websockets.go
package main

import (
	"encoding/json"
	"net/http"
)

// Message is the top-level object passed between the server and clients
type Message struct {
	MessageType string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
}

// ChatMessage is sent between clients via the server
type ChatMessage struct {
	Body      string `json:"body"`
	Username  string `json:"username"`
	TimeStamp string `json:"timeStamp"`
}

// UsernameMessage is sent by a client to register its username
type UsernameMessage struct {
	Body string `json:"body"`
}

// UsernameTakenMessage is sent by the server to a client
// in response to a UsernameMessage which tried to register a username which is in user
type UsernameTakenMessage struct {
	Body string `json:"body"`
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	nexus := newNexus()
	nexus.run()

	http.ListenAndServe(":8080", nil)
}
