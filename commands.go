package main

import (
	"bufio"
	"os"
	"strings"
)

type Action func([]byte, *ClientConnection)

var commandMap map[string]Action = map[string]Action{
	"login":  Login,
	"logout": Logout,
	"list":   UserList,
	"help":   Help,
}

func Help(message []byte, client *ClientConnection) {
	helpfile, err := os.Open("help.txt")

	if err != nil {
		client.Send(SystemMessageInfo("Impossible to find help file in server"))
		return
	}

	defer helpfile.Close()

	scanner := bufio.NewScanner(helpfile)
	for scanner.Scan() {
		client.Send(SystemMessageInfo(scanner.Text()))
	}
}

func Logout(message []byte, client *ClientConnection) {
	if client.logged {
		broadcastChannel <- SystemMessageLogout(client.name)
		client.Send(SystemMessageLoggedout(client.name))
		client.logged = false
	}
}

func Login(message []byte, client *ClientConnection) {
	newname := strings.TrimSpace(string(message))

	if len(newname) == 0 {
		client.Send(SystemMessageInfo("Please provide a name with the login command"))
		return
	}

	for _, otherclient := range connectionPool {
		if otherclient != nil && otherclient.logged && otherclient.name == newname {
			client.Send(SystemMessageInfo("That name is already selected!"))
			return
		}
	}

	if len(newname) > 16 {
		client.Send(SystemMessageInfo("Please use a shorter name"))
		return
	}

	if strings.IndexByte(newname, ':') >= 0 {
		client.Send(SystemMessageInfo("Please donde use colon in the name!"))
		return
	}

	if client.logged {
		oldname := client.name
		broadcastChannel <- SystemMessageRename(oldname, newname)
	} else {
		broadcastChannel <- SystemMessageLogin(newname)
	}

	client.Send(SystemMessageLogged(newname))
	client.name = newname
	client.logged = true
}

func UserList(message []byte, client *ClientConnection) {
	users := []string{}

	if client.logged {
		for _, otherclient := range connectionPool {
			if otherclient != nil && otherclient.logged {
				users = append(users, otherclient.name)
			}
		}
		client.Send(SystemMessageUserlist(users))
	} else {
		client.Send(SystemMessageInfo("You have to log in first"))
	}
}
