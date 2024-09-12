package main

type GameLobby interface {
	AddPlayer(userID string, score int) int
	PopTwoPlayers() (string, string)
}

type DefaultGameLobby struct {
	users map[string]int
}

func (gl *DefaultGameLobby) AddPlayer(userID string, score int) int {
	if _, ok := gl.users[userID]; !ok {
		gl.users[userID] = score
	}
	return len(gl.users)
}

func (gl *DefaultGameLobby) PopTwoPlayers() (string, string) {
	if len(gl.users) < 2 {
		return "", ""
	}

	var player1, player2 string
	count := 0

	for user := range gl.users {
		if count == 0 {
			player1 = user
		} else if count == 1 {
			player2 = user
		} else {
			break
		}
		count++
	}

	delete(gl.users, player1)
	delete(gl.users, player2)

	return player1, player2
}

func NewGameLobby() GameLobby {
	return &DefaultGameLobby{
		users: make(map[string]int),
	}
}
