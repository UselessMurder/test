package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/murder/chat/sessions"
	"log"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	intr, _ = c.Get("currentSession")
	currentSession = intr.(*sessions.Session)
	log.Println(currentSession)
	c.HTML(http.StatusOK, "index", nil)
	return
}
