package pool

import (
	"fmt"
	"log"
	"sort"
)

type Pool struct {
	inSession   bool               // indicator for game session
	Register    chan *Client       // for new user joining
	Unregister  chan *Client       // when a user quits the game
	Ping        chan Message       // for individual user message
	Broadcast   chan Message       // sending information to all clients
	Clients     map[string]*Client // map of client ID to Client object
	PlayerRoom  map[string]*Client // room for clients playing the game
	WaitingRoom map[string]*Client // waiting room ro users when game is in session
}

func NewPool() *Pool {
	return &Pool{
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Ping:        make(chan Message),
		Broadcast:   make(chan Message),
		Clients:     make(map[string]*Client),
		PlayerRoom:  make(map[string]*Client),
		WaitingRoom: make(map[string]*Client),
	}
}

//TODO: Remove all theses globals, better still encapsulate them into
// logical space.
var rs RankingSystem

func (p *Pool) Start() {
	for {
		if p.inSession && rs.randomGeneratorCount == 0 {
			go rs.broadCastRandomNumber(p)
		}

		select {
		case newClient := <-p.Register:
			// Inform all clients
			for _, client := range p.Clients {
				client.mu.Lock()
				if err := client.Conn.WriteJSON(Message{Type: "game-info", Info: "New User Joined..."}); err != nil {
					log.Println("pool: could not send JSON data to client")
				}
				client.mu.Unlock()
			}
			p.Clients[newClient.ID] = newClient

			// Give a response to the client alone
			newClient.mu.Lock()
			if err := newClient.Conn.WriteJSON(Message{Type: "game-info", Body: Body{ClientID: newClient.ID}, Info: "Welcome!"}); err != nil {
				log.Println("pool: could not send JSON data to client")
			}
			newClient.mu.Unlock()

			log.Println("Size of Connection Pool: ", len(p.Clients))
			break
		case client := <-p.Unregister:
			delete(p.Clients, client.ID)
			log.Println("Size of Connection Pool: ", len(p.Clients))
			for _, client := range p.Clients {
				client.mu.Lock()
				if err := client.Conn.WriteJSON(Message{Type: "game-info", Info: "User Disconnected..."}); err != nil {
					log.Println("pool: could not send JSON data to client")
				}
				client.mu.Unlock()
			}
			break
		case message := <-p.Ping:
			// update the client information
			if message.ClientID == "" {
				log.Println("client: message missing clientID")
				break
			}
			client := p.Clients[message.ClientID]
			if client.Name == "" {
				client.Name = message.Player
			}
			// when user wants to play
			if message.PlayerMode == "play" {
				client.Player = true

				if p.inSession { // if game is in session add to waiting room
					p.WaitingRoom[client.ID] = client
					log.Printf("game: %d players in waiting room\n", len(p.WaitingRoom))
					if err := client.Conn.WriteJSON(Message{Type: "player-wait", Info: "you can play when next game begins!"}); err != nil {
						log.Println("pool: could not send JSON data to client")
					}
				} else {
					// add the client to the playing room
					p.PlayerRoom[client.ID] = client
					log.Printf("game: %d players in game room\n", len(p.PlayerRoom))
				}
			}
			// for each round play
			if message.PlayerMode == "roundPlay" {
				// put the score in a sorted array
				roundInputs := []uint8{message.Input1, message.Input2}
				sort.Slice(roundInputs, func(i, j int) bool {
					return roundInputs[i] < roundInputs[j]
				})

				client.lowerBound = roundInputs[0]
				client.upperBound = roundInputs[1]
				break
			}
			break
		case message := <-p.Broadcast:
			log.Println("Sending message to all clients in Pool")
			for _, client := range p.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		}

		// if we have 2 or more players, start the game
		if len(p.PlayerRoom) >= 2 && !p.inSession {
			for _, client := range p.Clients {
				if err := client.Conn.WriteJSON(Message{Type: "game-start", Info: "game started!"}); err != nil {
					log.Println("pool: could not send JSON data to client")
				}
			}
			p.StartGame()
		}
	}
}

// StartGame will change game to be in session and notify all users.
func (p *Pool) StartGame() {
	p.inSession = true
	log.Println("game: Game is in session")
	for _, client := range p.Clients {
		if err := client.Conn.WriteJSON(Message{Type: "game-info", Info: "game in session"}); err != nil {
			log.Println("pool: could not send JSON data to client")
		}
	}
}
