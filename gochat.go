package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
)

var upgrader = websocket.Upgrader{}
var broadcastChannel chan (Serializable)

func broadcaster() {
	for {
		msg := <-broadcastChannel
		for _, client := range connectionPool {
			if client != nil && msg.BroadcastOk(client) {
				client.Send(msg)
			}
		}
	}
}

func chatConnectionHandler(w http.ResponseWriter, r *http.Request) {
	wsconnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	go startClientConnection(wsconnection)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Print("Running gochat on port ", port)

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/gochat", chatConnectionHandler)

	broadcastChannel = make(chan Serializable)
	go broadcaster()

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
