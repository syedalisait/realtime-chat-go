package server

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn            *websocket.Conn
	websocketServer *Websocket
	send            chan []byte
}

func NewClient(conn *websocket.Conn, websocketServer *Websocket) *Client {
	return &Client{
		conn:            conn,
		send:            make(chan []byte),
		websocketServer: websocketServer,
	}
}

func (c *Client) ReadMessages() {
	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			c.websocketServer.unregister <- c

			break
		}

		log.Println("chat message:", string(messageBytes))

		c.websocketServer.broadcast <- messageBytes
	}
}

func (c *Client) WriteMessages() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("cannot close websocket message: %s", err)
				}

				return
			}

			nextWriter, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("cannot create next writer: %s", err)

				return
			}

			_, err = nextWriter.Write(message)
			if err != nil {
				log.Printf("cannot write message: %s", err)

				return
			}

			if err = nextWriter.Close(); err != nil {
				log.Printf("cannot close writer: %s", err)

				return
			}
		}
	}
}
