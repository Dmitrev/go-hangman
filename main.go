package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"math/rand"
	"time"
	"os"
	"encoding/json"
)

type Word struct {
	correct []string
	guessed []string
}

type wordsSlice struct{
	Words []string `json:"words"`
}

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var word Word
var words wordsSlice
var upgrader websocket.Upgrader
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message) // broadcast channel
var lives int


func main() {
	// Set seeder for rand
	rand.Seed(time.Now().Unix())



	loadWords()
	generateNewWord()

	// Start the go routine here because server call is blocking
	go handleLetters()
	// Create simple file serverinfo info wa to serve static files
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)



	// Create websocket entry point
	http.HandleFunc("/ws", handleWebSocket)
	err := http.ListenAndServe(":8000", nil)

	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}

	log.Printf("Started server on port 8000")

}

func sendUpdate(ws *websocket.Conn) {
	ws.WriteJSON(&struct {
		Word []string `json:"word"`
		Lives int `json:"lives"`
	}{
		word.guessed,
		lives,
	})
}

func updateAll() {
	for client := range clients {
		sendUpdate(client)
	}
}

func generateNewWord() {
	// Reset lives
	lives = 10

	newWord := words.Words[rand.Intn(len(words.Words))]
	newWordSlice := strings.Split(newWord, "")
	emptySlice := make([]string, len(newWord))

	fmt.Println(newWord)
	// Create new word
	word = Word{
		correct: newWordSlice,
		guessed: emptySlice,
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
	}

	// Send current word to new connected client
	sendUpdate(ws)

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

		if msg.Type == "turn" {
			broadcast <- msg
		}


	}

	fmt.Println("user disconnected")

}


func handleLetters() {

	fmt.Println("handle letters")

	for {
		msg := <-broadcast

		fmt.Printf("Check letter %v", msg.Data)
		letter := msg.Data
		// Check if the letter is valid

		var letterIndexes []int

		// Check on which index this letter is
		for i, l := range word.correct {

			if strings.ToLower(letter) == strings.ToLower(l) {
				fmt.Printf("Found letter at index: %v", i)
				letterIndexes = append(letterIndexes, i)
			}
		}

		// Add to the guessed slice
		for _, li := range letterIndexes {
			fmt.Printf("Replacing letter at index: %v", li)
			word.guessed[li] = letter
		}

		if len(letterIndexes) == 0 {
			lives--
		}

		// send update
		for client := range clients {
			sendUpdate(client)
		}

		// Check if lost
		if lives == 0 {
			fmt.Printf("Game is over")
			generateNewWord()
			updateAll()
			return
		}

		// Check if game is over
		if strings.ToLower(strings.Join(word.guessed, "")) == strings.ToLower(strings.Join(word.correct, "")) {
			fmt.Printf("Game is over")
			generateNewWord()
			updateAll()
		}

	}
}

func loadWords(){
	f, err := os.Open("words.json")
	if err != nil {
		log.Fatalf("Failed to open file config")
	}

	dec := json.NewDecoder(f)
	err = dec.Decode(&words)
	f.Close()
	if err != nil {
		log.Fatalf("Bad JSON")
	}
}
