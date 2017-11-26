package routes

import (
	"../dbwrapper"
	"../models"
	"../sessions"
	_ "database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type tests_content struct {
	Page_status  int
	Requirements map[int]*tests_req
}

type tests_req struct {
	Wording  string
	Complete bool
}

func GetTestsHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	if !currentSession.Fixed {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	intr, _ = c.Get("database")
	db := intr.(*dbwrapper.DataBaseWrapper)
	tc := tests_content{2, make(map[int]*tests_req)}
	for k := range currentSession.Tests {
		var Req models.Requirement
		row, err := db.QueryRow("GetRequrinmentById", k)
		if err != nil {
			log.Fatal(err)
		}
		err = row.Scan(&Req.Id, &Req.Wording, &Req.Solution)
		if err != nil {
			log.Fatal(err)
		}
		tc.Requirements[k] = &tests_req{Req.Wording, currentSession.Tests[k].Complete}
	}
	c.HTML(http.StatusOK, "tests", tc)
}
