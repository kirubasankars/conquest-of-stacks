package main

type ScoreBoard struct {
	Player0 int `json:"player0"`
	Player1 int `json:"player1"`
}

type JoinGameMessage struct {
	Channel string
	Data    JoinGameData
}

type JoinGameData struct {
	GameCounter int    `json:"game_counter"`
	Action      string `json:"action"`
	GameID      string `json:"game_id"`

	Players []string `json:"players"`

	Score ScoreBoard `json:"score"`
}

type SegmentOccupiedData struct {
	GameCounter int    `json:"game_counter"`
	Action      string `json:"action"`
	GameID      string `json:"game_id"`

	SegmentID string     `json:"segment_id"`
	PlayerID  string     `json:"player_id"`
	Score     ScoreBoard `json:"score"`
}

type SegmentOccupiedMessage struct {
	Channel string
	Data    SegmentOccupiedData
}

type DefaultMessageData struct {
	GameCounter int        `json:"game_counter"`
	Action      string     `json:"action"`
	GameID      string     `json:"game_id"`
	UserID      string     `json:"user_id"`
	Score       ScoreBoard `json:"score"`
}

type DefaultMessage struct {
	Channel string             `json:"channel"`
	Data    DefaultMessageData `json:"data"`
}

type DiceRolledData struct {
	GameCounter int        `json:"game_counter"`
	Action      string     `json:"action"`
	GameID      string     `json:"game_id"`
	UserID      string     `json:"user_id"`
	Dices       []int      `json:"dices"`
	Sums        []int      `json:"sums"`
	Score       ScoreBoard `json:"score"`
}

type DiceRolledMessage struct {
	Channel string         `json:"channel"`
	Data    DiceRolledData `json:"data"`
}

type GameOverData struct {
	GameCounter int    `json:"game_counter"`
	GameID      string `json:"game_id"`
	Action      string `json:"action"`
}

type GameOverDataMessage struct {
	Channel string       `json:"channel"`
	Data    GameOverData `json:"data"`
}
