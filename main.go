package main

import (
	"./dbwrapper"
	"./routes"
	"./sessions"
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
	"log"
	"os/exec"
	"runtime"
	"time"
)

var sm sessions.SessionManager

func Middle(c *gin.Context) {
	log.Println("Connected", c.ClientIP())
	id := sm.GetCookie(c.Request, c.Writer)
	err, currentSession := sm.GetSession(id)
	if err != nil {
		currentSession = sessions.CreateSession(id, time.Now().Add(24*time.Hour))
		sm.SetSession(currentSession)
	}
	c.Set("currentSession", currentSession)
	c.Set("database", &dbwrapper.Wrapper)
	c.Next()
	sm.SetSession(currentSession)
	log.Println("Disconnected", c.ClientIP())
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func main() {
	log.Println("Start listening 8080")
	err := dbwrapper.Wrapper.ReplaceRequestList("requests.sqls")
	if err != nil {
		log.Panicln("Sql error:", err)
	}
	sm.OpenSessionManager()
	defer sm.CloseSessionManager()
	r := gin.Default()
	r.Use(favicon.New("./favicon.ico"))
	r.Use(Middle)
	r.GET("/", routes.GetIndexHandler)
	r.GET("/tests", routes.GetTestsHandler)
	r.GET("/tests/questions", routes.GetQuestionsHandler)
	r.GET("/results", routes.GetResultsHandler)
	r.GET("/results/report", routes.GetReportHandler)
	r.POST("/", routes.PostIndexHandler)
	r.POST("/tests/questions", routes.PostQuestionsHandler)
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")
	go open("http://localhost:8080/")
	r.Run(":8080")
}
