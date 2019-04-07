// websockets.go
package main

import (
	"encoding/json"
	"net/http"
)

type ChatMessage struct {
	Body      string `json:"body"`
	Username  string `json:"username"`
	TimeStamp string `json:"timeStamp"`
}

type UsernameMessage struct {
	Body string `json:"body"`
}

type UsernameTakenMessage struct {
	Body string `json:"body"`
}

type Message struct {
	MessageType string          `json:"type"`
	Payload     json.RawMessage `json:"payload"`
}

func main() {
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	nexus := newNexus()
	nexus.run()

	http.ListenAndServe(":8080", nil)
}
