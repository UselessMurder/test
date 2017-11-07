package main

import (
	"./routes"
	"./sessions"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var sm sessions.SessionManager

func Middle(c *gin.Context) {
	log.Println("Connected", c.ClientIP())
	id := userSessions.GetCookie(c.Request, c.Writer)
	err, currentSession := userSessions.GetSession(id)
	if err != nil {
		currentSession = &sessions.Session{id, false, 0, 0, time.Now().Add(24 * time.Hour)}
		userSessions.SetSession(currentSession)
	}
	c.Set("currentSession", currentSession)
	c.Next()
	userSessions.SetSession(currentSession)
	log.Println("Disconnected", c.ClientIP())
}

func main() {
	log.Println("Start listening 8080")
	r := gin.Default()
	r.Use(Middle)
	r.GET("/", routes.IndexHandler)
	r = gin.LoadHTMLGlob("templates/*")
	r = gin.Static("/assets", "./assets")
	r.Run(":8080")
}
