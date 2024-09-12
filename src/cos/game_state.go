package main

type GameState struct {
	MutationCounter int `json:"game_counter"`

	GameID  string   `json:"gameID"`
	Players []string `json:"players"`
	Stacks  map[string]*GameSegment

	TimeOutInSecs int    `json:"timeout_in_secs"`
	CurrentPlayer int    `json:"currentPlayer"`
	CurrentAction string `json:"currentAction"`

	Player0 []string `json:"player_0"`
	Player1 []string `json:"player_1"`
}

func (gameState *GameState) GetCurrentUserID() string {
	return gameState.Players[gameState.CurrentPlayer]
}

func (gameState *GameState) GetOtherUserID() string {
	if gameState.CurrentPlayer == 0 {
		return gameState.Players[1]
	}
	return gameState.Players[0]
}
