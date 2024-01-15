package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	//this here upgrades http connection to websocket connection
	//CheckOrigin allows access from every origin basicallty [*]
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

//`` are used to specify json fieldnames, used while marshalling and unmarshalling packets

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to chat room")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	clients[conn] = true

	for {
		var msg Message

		err := conn.ReadJSON(&msg)

		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)

			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {

	http.HandleFunc("/", homePage)
    http.HandleFunc("/ws", handleConnections)

	go handleMessages()

	fmt.Println("Server started on :8000")
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		panic("Error starting on server" + err.Error())
	}
}
