package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var gameEngine *GameEngine

func main() {
	gameEngine = NewGameEngine()
	router := gin.Default()
	router.Static("/www", "./www")

	users["k"] = "k"
	users["s"] = "s"

	jwtMiddleware, err := jwt.New(jwtParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	router.Use(handlerMiddleWare(jwtMiddleware))
	registerRoute(router, jwtMiddleware)

	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
