package main

import (
	"encoding/json"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func joinGameHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	userID, ok := claims[identityKey].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	secs := c.Param("secs")
	gameEngine.PutUserInGameIfPossible(userID, secs)
	c.Status(http.StatusOK)
}

func connectedGameHandler(c *gin.Context) {
	gameID := c.Param("game_id")
	claims := jwt.ExtractClaims(c)
	userID, ok := claims[identityKey].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	game, exists := gameEngine.games[gameID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	game.Lock()
	defer game.Unlock()

	game.activePlayers[userID] = nil

	if len(game.activePlayers) == 2 {
		gameEngine.StartGame(gameID)
	}
	c.Status(http.StatusOK)
}

func occupySegmentGameHandler(c *gin.Context) {
	gameID := c.Param("game_id")
	segmentID := c.Param("segment_id")
	claims := jwt.ExtractClaims(c)
	userID, ok := claims[identityKey].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	_, exists := gameEngine.games[gameID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	gameEngine.OccupySegment(userID, gameID, segmentID)
	c.Status(http.StatusOK)
}

func rollGameHandler(c *gin.Context) {
	gameID := c.Param("game_id")
	claims := jwt.ExtractClaims(c)
	userID, ok := claims[identityKey].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	_, exists := gameEngine.games[gameID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	gameEngine.RollDice(userID, gameID)
	c.Status(http.StatusOK)
}

func endTurnGameHandler(c *gin.Context) {
	gameID := c.Param("game_id")
	claims := jwt.ExtractClaims(c)
	userID, ok := claims[identityKey].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token claims"})
		return
	}

	_, exists := gameEngine.games[gameID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	gameEngine.EndTurn(gameID, userID)

	c.Status(http.StatusOK)
}

func getGameStateHandler(c *gin.Context) {
	var board GameBoard
	byt, err := os.ReadFile("www/board.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read board file"})
		return
	}

	err = json.Unmarshal(byt, &board)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal board"})
		return
	}

	gameID := c.Param("game_id")
	game, exists := gameEngine.games[gameID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	board.State = *game
	c.JSON(http.StatusOK, board)
}
