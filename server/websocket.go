package server

import "log"

type Websocket struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	unregister chan *Client
}

func NewWebsocket() *Websocket {
	return &Websocket{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (ws *Websocket) Initialize() {
	for {
		select {
		case client := <-ws.Register:
			ws.clients[client] = true
			log.Println("registered client successfully")
		case client := <-ws.unregister:
			delete(ws.clients, client)
			log.Println("unregistered client successfully")
		case message := <-ws.broadcast:
			log.Println("publishing message to other clients", string(message))
			for client, isActive := range ws.clients {
				if isActive {
					client.send <- message
				}
			}
		}
	}
}
