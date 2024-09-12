package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Game struct {
	mu                   *sync.Mutex
	ticker               *time.Ticker
	tickerDone           chan bool
	lastRolledSums       []int
	currentActionTimeOut time.Time
	activePlayers        map[string]interface{}
	GameState
	communicator *Communicator
}

func (game *Game) Lock() {
	game.mu.Lock()
}

func (game *Game) Unlock() {
	game.mu.Unlock()
}

func (game *Game) ChooseSegment() {
	game.MutationCounter++

	chooseSegment := DefaultMessage{
		Channel: game.GameID,
		Data: DefaultMessageData{
			GameCounter: game.MutationCounter,
			Action:      "choose_segment",
			GameID:      game.GameID,
			UserID:      game.GetCurrentUserID(),
			Score:       game.GetSourceBoard(),
		},
	}

	byt, _ := json.Marshal(chooseSegment)
	game.communicator.Sent(string(byt))
}

func remove(slice []string, item string) []string {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func find(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func (game *Game) OwnSegment(segment string) {
	if segment == "_none_" {
		game.ScheduleNextAction(ActionEndTurn)
		game.currentActionTimeOut = time.Now()
		return
	}

	segment = strings.TrimSpace(segment)
	if segment == "" || !strings.Contains(segment, "_") || game.CurrentAction != ActionChooseSegment {
		return
	}

	items := strings.Split(segment, "_")
	stackID := items[0]
	stack := game.Stacks[stackID]

	if !find(game.lastRolledSums, stack.ControlGroup) {
		return
	}

	stack.Segments = remove(stack.Segments, segment)

	game.MutationCounter++
	if game.CurrentPlayer == 0 {
		game.Player0 = append(game.Player0, segment)
		stack.Player0 = append(stack.Player0, segment)
	} else {
		game.Player1 = append(game.Player1, segment)
		stack.Player1 = append(stack.Player1, segment)
	}

	currentPlayerID := game.GetCurrentUserID()
	segmentOccupiedMessage := SegmentOccupiedMessage{
		Channel: game.GameID,
		Data: SegmentOccupiedData{
			GameCounter: game.MutationCounter,
			Action:      "segment_occupied",
			GameID:      game.GameID,
			SegmentID:   segment,
			PlayerID:    currentPlayerID,
			Score:       game.GetSourceBoard(),
		},
	}

	byt, _ := json.Marshal(segmentOccupiedMessage)
	game.communicator.Sent(string(byt))
	game.ScheduleNextAction(ActionEndTurn)
	game.currentActionTimeOut = time.Now()

	game.CheckGameOverAndAnnounce()
}

func (game *Game) GetSourceBoard() ScoreBoard {
	player0Score := 0
	player1Score := 0

	for _, segment := range game.Stacks {
		if len(segment.Player0) > 0 && math.Ceil(float64(segment.NumberOfSegments)/2) <= float64(len(segment.Player0)) {
			player0Score += segment.Points
		}
		if len(segment.Player1) > 0 && math.Ceil(float64(segment.NumberOfSegments)/2) <= float64(len(segment.Player1)) {
			player1Score += segment.Points
		}
	}

	return ScoreBoard{player0Score, player1Score}
}

func (game *Game) EndTurn() {
	game.MutationCounter++
	currentPlayerID := game.GetCurrentUserID()

	game.CurrentPlayer = 1 - game.CurrentPlayer

	endTurn := DefaultMessage{
		Channel: game.GameID,
		Data: DefaultMessageData{
			GameCounter: game.MutationCounter,
			Action:      "end_turn",
			GameID:      game.GameID,
			UserID:      currentPlayerID,
			Score:       game.GetSourceBoard(),
		},
	}

	byt, _ := json.Marshal(endTurn)
	game.communicator.Sent(string(byt))
}

func (game *Game) StartRoll() {
	currentPlayerID := game.GetCurrentUserID()
	game.MutationCounter++

	startRoll := DefaultMessage{
		Channel: game.GameID,
		Data: DefaultMessageData{
			GameCounter: game.MutationCounter,
			Action:      "start_roll",
			GameID:      game.GameID,
			UserID:      currentPlayerID,
			Score:       game.GetSourceBoard(),
		},
	}

	byt, _ := json.Marshal(startRoll)
	game.communicator.Sent(string(byt))
}

func (game *Game) RollDice() {
	game.MutationCounter++
	var nums []int
	var uniqSums []int

	for i := 0; i < 4; i++ {
		nums = append(nums, rand.Intn(8)+3) // Ensuring dice rolls between 3 and 10
	}
	sums := []int{
		nums[0] + nums[1],
		nums[0] + nums[2],
		nums[0] + nums[3],
		nums[1] + nums[2],
		nums[1] + nums[3],
		nums[2] + nums[3],
	}

	for _, x := range sums {
		if !find(uniqSums, x) {
			uniqSums = append(uniqSums, x)
		}
	}

	//uniqSums = []int{7, 11, 14}

	diceRolledMessage := DiceRolledMessage{
		Channel: game.GameID,
		Data: DiceRolledData{
			GameCounter: game.MutationCounter,
			Action:      "dice_rolled",
			GameID:      game.GameID,
			UserID:      game.GetCurrentUserID(),
			Sums:        uniqSums,
			Dices:       nums,
			Score:       game.GetSourceBoard(),
		},
	}

	game.lastRolledSums = uniqSums
	byt, _ := json.Marshal(diceRolledMessage)
	game.communicator.Sent(string(byt))

	game.ChooseSegment()
}

func (game *Game) run() {
	game.Lock()
	defer game.Unlock()

	if !game.IsActionTimeout() {
		return
	}

	switch game.CurrentAction {
	case ActionRollDice:
		game.RollDice()
		game.ScheduleNextAction(ActionChooseSegment)
	case ActionChooseSegment:
		game.OwnSegment(game.getMajoritySegment(game.CurrentPlayer))
	case ActionEndTurn:
		game.EndTurn()
		game.StartRoll()
		game.ScheduleNextAction(ActionRollDice)
	}
}

func (game *Game) Start() {
	if game.CurrentAction != "" {
		return
	}
	game.StartRoll()
	game.ScheduleNextAction(ActionRollDice)

	go func() {
		for {
			select {
			case <-game.tickerDone:
				return
			case <-game.ticker.C:
				game.run()
			}
		}
	}()

	currentPlayerID := game.GetCurrentUserID()
	startRoll := DefaultMessage{
		Channel: game.GameID,
		Data: DefaultMessageData{
			GameCounter: game.MutationCounter,
			Action:      "game_started",
			GameID:      game.GameID,
			UserID:      currentPlayerID,
			Score:       game.GetSourceBoard(),
		},
	}

	byt, _ := json.Marshal(startRoll)
	game.communicator.Sent(string(byt))
}

func (game *Game) Stop() {
	game.tickerDone <- true
	game.ticker.Stop()
}

func (game *Game) ScheduleNextAction(action string) {
	game.CurrentAction = action
	game.currentActionTimeOut = time.Now().Add(time.Duration(game.TimeOutInSecs+2) * time.Second)

	if game.CurrentAction == ActionEndTurn {
		game.currentActionTimeOut = time.Now().Add(5 * time.Second)
	}
}

func (game *Game) IsActionTimeout() bool {
	return time.Now().After(game.currentActionTimeOut)
}

func (game *Game) getMajoritySegment(currentPlayer int) string {
	for _, sum := range game.lastRolledSums {
		for _, stack := range game.Stacks {
			if stack.ControlGroup == sum && hasMajority(stack, currentPlayer) {
				if len(stack.Segments) > 0 {
					return stack.Segments[0]
				}
			}
		}
	}
	return "_none_"
}

func (game *Game) CheckGameOverAndAnnounce() {
	gameOver := true
	for _, stack := range game.Stacks {
		if len(stack.Segments) > 0 {
			gameOver = false
		}
	}

	if gameOver {
		gameOverMessage := GameOverDataMessage{
			Channel: game.GameID,
			Data: GameOverData{
				GameCounter: game.MutationCounter,
				Action:      "game_over",
				GameID:      game.GameID,
			},
		}
		byt, _ := json.Marshal(gameOverMessage)
		game.communicator.Sent(string(byt))
		game.Stop()
	}
}

func hasMajority(stack *GameSegment, currentPlayer int) bool {
	count := len(stack.Segments)
	total := 0
	if currentPlayer == 0 {
		total = len(stack.Player0)
	} else {
		total = len(stack.Player1)
	}
	return count > total/2
}

func NewGame(gameID string, timeOutInSecs int, players []string) *Game {
	var board GameBoard
	byt, err := os.ReadFile("www/board.json")
	if err != nil {
		panic(fmt.Sprintf("failed to read board file: %v", err))
	}
	if err := json.Unmarshal(byt, &board); err != nil {
		panic(fmt.Sprintf("failed to unmarshal board: %v", err))
	}

	stacks := make(map[string]*GameSegment)
	for _, segment := range board.Board {
		stack := &GameSegment{
			ID:               segment.ID,
			NumberOfSegments: segment.NumberOfSegments,
			Points:           segment.Points,
			ControlGroup:     segment.ControlGroup,
		}
		for i := 0; i < segment.NumberOfSegments; i++ {
			stack.Segments = append(stack.Segments, fmt.Sprintf("%s_%d", segment.ID, i))
		}
		stacks[segment.ID] = stack
	}

	return &Game{
		GameState: GameState{
			MutationCounter: 0,
			GameID:          gameID,
			Players:         players,
			CurrentPlayer:   0,
			CurrentAction:   "",
			Stacks:          stacks,
			TimeOutInSecs:   timeOutInSecs,
		},
		mu:            &sync.Mutex{},
		ticker:        time.NewTicker(1 * time.Second),
		tickerDone:    make(chan bool),
		activePlayers: make(map[string]interface{}),
		communicator:  NewCommunicator(),
	}
}
