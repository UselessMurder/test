package routes

import (
	"../sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func IndexHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	log.Println(currentSession)
	c.HTML(http.StatusOK, "index", nil)
	return
}
