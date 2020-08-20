package pool

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type ScoreBoard []*Client

// The following 3 methods are needed to implement the sort.Interface
// which will be used to rank the players on ScoreBoard.
//
// It uses the following criteria in order of priority:
// TotalScore, upperBound, lowerBound and name.
func (s ScoreBoard) Len() int { return len(s) }
func (s ScoreBoard) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ScoreBoard) Less(i, j int) bool {
	if s[i].TotalScore == s[j].TotalScore {
		// check the upperbound
		if s[i].upperBound == s[j].upperBound {
			// check lowerbound
			if s[i].lowerBound == s[j].lowerBound {
				// check their name
				return s[i].Name < s[j].Name
			}
			return s[i].lowerBound > s[j].lowerBound
		}
		return s[i].upperBound > s[j].upperBound
	}
	return s[i].TotalScore > s[j].TotalScore
}

// rnaking system
type RankingSystem struct {
	scoreBoard           ScoreBoard
	round                uint8
	randomGeneratorCount uint8
}

// searchClient will search for the position of a client in the scoreboard,
// and return their position and -1 if not found.
func (rs *RankingSystem) searchClient(c *Client) *Client {
	for leftPointer, rightPointer := 0, len(rs.scoreBoard)-1; leftPointer < rightPointer; leftPointer, rightPointer = leftPointer+1, rightPointer-1 {
		if found := rs.scoreBoard[leftPointer]; found.ID == c.ID {
			return found
		} else if found := rs.scoreBoard[rightPointer]; found.ID == c.ID {
			return found
		}
	}

	return nil
}

func (rs *RankingSystem) addClientToScoreBoard(c *Client) []*Client {
	// search for the client
	client := rs.searchClient(c)
	if client != nil { // if client already exist
		client.TotalScore = c.TotalScore
	} else {
		rs.scoreBoard = append(rs.scoreBoard, c)
	}
	sort.Sort(rs.scoreBoard)

	return rs.scoreBoard
}

func (rs *RankingSystem) generateRandomNumber() uint8 {
	return uint8(rand.Intn(11-1) + 1)
}

func (rs *RankingSystem) ResetGame(p *Pool) {
	rs.round = 0
	p.inSession = false
	// empty the playroom
	p.PlayerRoom = make(map[string]*Client)
	fmt.Printf("game: players room emptied: %d\n", len(p.WaitingRoom))

	// wait for 15 secs before game restarts
	p.Broadcast <- Message{Type: "game-info", Info: fmt.Sprintf("new-game: starts in %d secs\n", 15)}
	time.Sleep(time.Second * 15)

	// add players from waiting room to player room
	for ID, client := range p.WaitingRoom {
		p.PlayerRoom[ID] = client
	}
	// empty the waiting room
	p.WaitingRoom = make(map[string]*Client)
}

func (rs *RankingSystem) broadCastRandomNumber(p *Pool) {
	for {
		if p.inSession {
			// wait 2 secs for each round
			time.Sleep(time.Second * 2)
			rs.round++
			rs.randomGeneratorCount++

			randNum := rs.generateRandomNumber()
			for _, client := range p.PlayerRoom {
				client.mu.Lock()

				// Exact Match: add 5 to total score
				if randNum == client.lowerBound || randNum == client.upperBound {
					client.TotalScore += 5
					client.Conn.WriteJSON(Message{Type: "play-info", Info: fmt.Sprintf("game: exact match!: %d", randNum)})
					// client.mu.Unlock()
				} else if randNum > client.lowerBound && randNum < client.upperBound {
					// client.mu.Lock()
					client.Conn.WriteJSON(Message{Type: "play-info", Info: fmt.Sprintf("game: Nice! you guessed right!: %d", randNum)})
					// client.mu.Unlock()
					// calculate the score and add client to the ranking system (priority queue)
					// Inside bounds:+5 - (upper bound - lower bound)
					// -1 wil be deducted if not within bounds
					client.TotalScore += int8((5 - (client.upperBound - client.lowerBound)))
				} else {
					// -1 wil be deducted if not within bounds
					// client.mu.Lock()
					client.TotalScore--
					client.Conn.WriteJSON(Message{Type: "play-info", Info: fmt.Sprintf("game: better luck next time!: %d", randNum)})
					// client.mu.Unlock()
				}
				client.mu.Unlock()

				// add user to the scoreboard
				rs.addClientToScoreBoard(client)

				// if any player reaches 21 score, end the game
				// and declare the player as the winner
				if client.TotalScore == 21 {
					fmt.Print("game: 21 points score, we have a winner! quitting...")
					rs.ResetGame(p)
					continue
				}
			}

			if rs.round > 30 {
				fmt.Print("game: reach 30 rounds, quitting...")
				// reset everything
				rs.ResetGame(p)
				continue
			}
			fmt.Printf("game: scoreboard after round %d: %v\n", rs.round, rs.scoreBoard)
			// broadcast ScoreBoard to all clients
			p.Broadcast <- Message{Type: "scoreboard", Info: fmt.Sprintf("game: scoreboard after round %d\n", rs.round), Body: Body{ScoreBoard: rs.scoreBoard}}
		}
	}
}
