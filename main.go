package main

import (
	"gin_example/controllers/auth"
	"gin_example/core"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.POST("/auth/signin/", auth.Sigin)
	// /login/

	return r
}

func main() {
	core.InitDB()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
