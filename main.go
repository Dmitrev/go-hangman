package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Word struct {
	correct []string
	guessed []string
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var word Word
var upgrader websocket.Upgrader
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message) // broadcast channel

func main() {
	// Create simple file serverinfo info wa to serve static files
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// Create new word
	word = Word{
		correct: []string{"H", "E", "L", "L", "O"},
		guessed: []string{"", "", "", "", ""},
	}

	// Create websocket entry point
	http.HandleFunc("/ws", handleWebSocket)
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
	go handleLetters()

	log.Printf("Started server on port 8000")

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
	}

	// Send current word to new connected client
	ws.WriteJSON(&struct {
		Word []string `json:"word"`
	}{
		word.guessed,
	})

	// Close set close connection
	defer ws.Close()

	var msg Message
	log.Printf("User connected")
	clients[ws] = true

	// ** GAME LOOP ** //
	for {
		err := ws.ReadJSON(&msg)

		// If we cannot read the message anymore that means that the user is disconnected
		if err != nil {
			fmt.Println("CLOSE")
			fmt.Println(err)
			delete(clients, ws)
			break
		}

		fmt.Println("Check")

		broadcast <- msg
	}

	fmt.Println("user disconnected")

}

func handleLetters() {
	for {
		msg := <-broadcast

		fmt.Printf("Check letter %v", msg.Data)
		// Check if the letter is valid

		// var letterIndexes []int
		//
		// // Check on which index this letter is
		// for i, l := range word.correct {
		// 	if letter == l {
		// 		letterIndexes = append(letterIndexes, i)
		// 	}
		// }
		//
		// // Add to the guessed slice
		// for _, li := range letterIndexes {
		// 	word.guessed[li] = letter
		// }

		// send update

	}
}
