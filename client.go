package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	nexus    Nexus
	send     chan ChatMessage
	conn     *websocket.Conn
	username string
}

func (c *Client) streamIn() {
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(c.nexus.clients, c)
			break
		}

		switch msg.MessageType {
		case "username":
			usernameMessage := UsernameMessage{}
			err = json.Unmarshal(msg.Payload, &usernameMessage)
			if err != nil {
				log.Printf("error: %v", err)
				delete(c.nexus.clients, c)
				break
			}
			var matched bool
			for client, _ := range c.nexus.clients {
				if client.username == usernameMessage.Body {
					matched = true
					break
				}
			}
			if !matched {
				c.username = usernameMessage.Body
				fmt.Println(c)
			} else {
				log.Printf("error: Username is already taken")
				errorMessage := UsernameTakenMessage{Body: "Username is already taken"}
				b, err := json.Marshal(errorMessage)
				if err != nil {
					log.Printf("error: %v", err)
					delete(c.nexus.clients, c)
					break
				}
				message := Message{Payload: b, MessageType: "error"}
				c.conn.WriteJSON(message)
			}

		case "chat":
			chatMessage := ChatMessage{}
			err = json.Unmarshal(msg.Payload, &chatMessage)
			if err != nil {
				log.Printf("error: %v", err)
				delete(c.nexus.clients, c)
				break
			}
			c.nexus.broadcast <- chatMessage
		}
	}

}

func (c *Client) streamOut() {
	for {
		msg := <-c.send
		log.Println("sending message " + string(msg.Body))

		err := c.conn.WriteJSON(msg)
		if err != nil {
			log.Printf("error sending: %v", err)
			c.conn.Close()
			delete(c.nexus.clients, c)
		}

	}

}
