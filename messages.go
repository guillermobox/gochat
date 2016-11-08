package main

import (
	"fmt"
	"strings"
)

type Serializable interface {
	Serialize() []byte
	BroadcastOk(*ClientConnection) bool
}

type ChatMessage struct {
	message string
	client  *ClientConnection
}

func (m ChatMessage) Serialize() []byte {
	return []byte("chat:" + m.client.name + ":" + m.message)
}

func (m ChatMessage) BroadcastOk(client *ClientConnection) bool {
	return client.logged && m.client.logged && (client != m.client)
}

type SystemMessage struct {
	msgtype string
	message string
}

func (m SystemMessage) Serialize() []byte {
	return []byte(fmt.Sprintf("system:%s:%s", m.msgtype, m.message))
}

func (m SystemMessage) BroadcastOk(client *ClientConnection) bool {
	return client.logged
}

func SystemMessageInfo(msg string) SystemMessage {
	return SystemMessage{"info", msg}
}

func SystemMessageLogged(name string) SystemMessage {
	return SystemMessage{"logged", name}
}

func SystemMessageLoggedout(name string) SystemMessage {
	return SystemMessage{"loggedout", name}
}

func SystemMessageRename(oldname, newname string) SystemMessage {
	return SystemMessage{"rename", oldname + ":" + newname}
}

func SystemMessageLogin(name string) SystemMessage {
	return SystemMessage{"login", name}
}

func SystemMessageLogout(name string) SystemMessage {
	return SystemMessage{"logout", name}
}

func SystemMessageUserlist(names []string) SystemMessage {
	return SystemMessage{"userlist", strings.Join(names, ":")}
}
