package main

import (
	"github.com/gorilla/websocket"
	"github.com/syedalisait/realtime-chat-go/server"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	websocketServer := server.NewWebsocket()
	go websocketServer.Initialize()

	initializeRouter(websocketServer)

	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func initializeRouter(websocketServer *server.Websocket) {
	fs := http.FileServer(http.Dir("./frontend/public"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocketHandler(websocketServer, w, r)
	})
}

func websocketHandler(websocketServer *server.Websocket, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)

		return
	}

	log.Println("creating a new client")

	client := server.NewClient(conn, websocketServer)

	go client.ReadMessages()
	go client.WriteMessages()

	websocketServer.Register <- client
}
