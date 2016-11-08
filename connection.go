package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"sync"
)

type ClientConnection struct {
	connection *websocket.Conn
	logged     bool
	name       string
	mutex      sync.Mutex
}

var connectionPool [16]*ClientConnection

func startClientConnection(wsconnection *websocket.Conn) {
	client := &ClientConnection{connection: wsconnection}

	for i, slot := range connectionPool {
		if slot == nil {
			connectionPool[i] = client
			break
		}
	}

	client.Loop()
}

func (client *ClientConnection) Loop() {
	defer client.Close()

	for {
		_, message, err := client.connection.ReadMessage()

		if err != nil {
			return
		}

		if len(message) == 0 {
			continue
		}

		client.Process(message)
	}
}

func (client *ClientConnection) Close() {
	client.connection.Close()
	for i, otherclient := range connectionPool {
		if otherclient == client {
			connectionPool[i] = nil
			break
		}
	}
	if client.logged {
		broadcastChannel <- SystemMessage{"logout", client.name}
	}
}

func (client *ClientConnection) Send(msg Serializable) {
	payload := msg.Serialize()
	client.mutex.Lock()
	defer client.mutex.Unlock()
	client.connection.WriteMessage(websocket.TextMessage, payload)
}

func (client *ClientConnection) Process(msg []byte) {
	start := bytes.IndexByte(msg, '/')
	end := bytes.IndexByte(msg, ' ')

	if start != 0 {
		if client.logged {
			broadcastChannel <- ChatMessage{string(msg), client}
		} else {
			client.Send(SystemMessageInfo("You have to log in first"))
		}
	} else {
		if end == -1 {
			end = len(msg)
		}
		command := string(msg[start+1 : end])
		payload := msg[end:]
		action, ok := commandMap[command]
		if ok {
			action(payload, client)
		} else {
			client.Send(SystemMessageInfo("Unknown command"))
		}
	}
}
