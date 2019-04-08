package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Client is a wrapper for a websocket connection
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
			err = c.handleUsernameMessage(msg.Payload)
			if err != nil {
				panic(err)
			}
		case "chat":
			err = c.handleChatMessage(msg.Payload)
			if err != nil {
				panic(err)
			}
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

func (c *Client) handleUsernameMessage(payload json.RawMessage) error {
	usernameMessage := UsernameMessage{}
	err := json.Unmarshal(payload, &usernameMessage)
	if err != nil {
		log.Printf("error: %v", err)
		delete(c.nexus.clients, c)
		return errors.New("Couldn't unmarshal msg.Payload into a usernameMessage")
	}
	var matched bool
	for client := range c.nexus.clients {
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
			return errors.New("Couldn't marshal errorMessage")
		}
		message := Message{Payload: b, MessageType: "error"}
		c.conn.WriteJSON(message)

	}
	return nil
}

func (c *Client) handleChatMessage(payload json.RawMessage) error {
	chatMessage := ChatMessage{}
	err := json.Unmarshal(payload, &chatMessage)
	if err != nil {
		log.Printf("error: %v", err)
		delete(c.nexus.clients, c)
		return errors.New("Couldn't unmarshal payload into chatMessage")
	}
	c.nexus.broadcast <- chatMessage
	return nil
}
