package main

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func registerRoute(r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	r.POST("/login", handle.LoginHandler)
	r.POST("/register", registerHandler)
	r.NoRoute(handle.MiddlewareFunc(), handleNoRoute())

	auth := r.Group("/game", handle.MiddlewareFunc())
	auth.GET("/refresh_token", handle.RefreshHandler)
	auth.POST("/join_game/:secs", joinGameHandler)

	auth.GET("/:game_id", getGameStateHandler)
	auth.POST("/:game_id/connected", connectedGameHandler)
	auth.POST("/:game_id/roll", rollGameHandler)
	auth.POST("/:game_id/end_turn", endTurnGameHandler)
	auth.POST("/:game_id/occupy/:segment_id", occupySegmentGameHandler)
}
