package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"sync"
)

type GameEngine struct {
	lobby15s GameLobby
	lobby30s GameLobby

	mu           sync.Mutex
	games        map[string]*Game
	communicator *Communicator
}

func (engine *GameEngine) PutUserInGameIfPossible(userID string, lobby string) {
	var count int
	engine.mu.Lock()
	defer engine.mu.Unlock()

	switch lobby {
	case "15s":
		count = engine.lobby15s.AddPlayer(userID, 0)
	case "30s":
		count = engine.lobby30s.AddPlayer(userID, 0)
	default:
		return
	}

	if count > 1 {
		gameID := uuid.New().String()
		var player1, player2 string
		if lobby == "15s" {
			player1, player2 = engine.lobby15s.PopTwoPlayers()
		} else {
			player1, player2 = engine.lobby30s.PopTwoPlayers()
		}
		game := NewGame(gameID, getTimeOut(lobby), []string{player1, player2})
		engine.JoinGame(game)
	}
}

func getTimeOut(lobby string) int {
	switch lobby {
	case "15s":
		return 15
	case "30s":
		return 30
	default:
		return 0
	}
}

func (engine *GameEngine) JoinGame(game *Game) {
	engine.games[game.GameID] = game

	joinGameMessage := JoinGameMessage{
		Channel: "lobby",
		Data: JoinGameData{
			GameCounter: game.MutationCounter,
			GameID:      game.GameID,
			Action:      "join_game",
			Players:     game.Players,
		},
	}
	byt, err := json.Marshal(joinGameMessage)
	if err == nil {
		engine.communicator.Sent(string(byt))
	}
}

func (engine *GameEngine) StartGame(gameID string) {
	game, exists := engine.games[gameID]
	if !exists {
		return
	}
	game.Start()
}

func (engine *GameEngine) RollDice(userID, gameID string) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	game, exists := engine.games[gameID]
	if !exists {
		return
	}

	game.Lock()
	defer game.Unlock()
	if game.GetCurrentUserID() != userID || game.CurrentAction != ActionRollDice {
		return
	}
	game.RollDice()
	game.ScheduleNextAction(ActionChooseSegment)
}

func (engine *GameEngine) EndTurn(gameID string, userID string) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	game, exists := engine.games[gameID]
	if !exists {
		return
	}

	game.Lock()
	defer game.Unlock()

	if game.GetCurrentUserID() != userID || game.CurrentAction != ActionChooseSegment {
		return
	}

	game.EndTurn()
	game.StartRoll()
	game.ScheduleNextAction(ActionRollDice)
}

func (engine *GameEngine) OccupySegment(userID, gameID, segmentID string) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	game, exists := engine.games[gameID]
	if !exists {
		return
	}

	game.Lock()
	defer game.Unlock()
	if userID != game.GetCurrentUserID() {
		return
	}

	game.OwnSegment(segmentID)
}

func NewGameEngine() *GameEngine {
	return &GameEngine{
		games:        make(map[string]*Game),
		lobby15s:     NewGameLobby(),
		lobby30s:     NewGameLobby(),
		communicator: NewCommunicator(),
	}
}
