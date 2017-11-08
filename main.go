package main

import (
	"./routes"
	"./sessions"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var sm sessions.SessionManager

func Middle(c *gin.Context) {
	log.Println("Connected", c.ClientIP())
	id := sm.GetCookie(c.Request, c.Writer)
	err, currentSession := sm.GetSession(id)
	if err != nil {
		currentSession = &sessions.Session{id, false, 0, 0, time.Now().Add(24 * time.Hour)}
		sm.SetSession(currentSession)
	}
	c.Set("currentSession", currentSession)
	c.Next()
	sm.SetSession(currentSession)
	log.Println("Disconnected", c.ClientIP())
}

func main() {
	log.Println("Start listening 8080")
	sm.OpenSessionManager()
	defer sm.CloseSessionManager()
	r := gin.Default()
	r.Use(Middle)
	r.GET("/", routes.IndexHandler)
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	r.Run(":8080")
}
