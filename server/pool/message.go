package pool

type Body struct {
	Player     string `json:"player"`
	ClientID   string `json:"clientID"`
	Input1     uint8  `json:"input1"`
	Input2     uint8  `json:"input2"`
	PlayerMode string `json:"playerMode"`
	ScoreBoard `json:"scoreboard"`
}

type Message struct {
	Type string `json:"type"`
	Info string `json:"info"`
	Body
}
