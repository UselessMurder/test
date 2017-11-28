package routes

import (
	"../dbwrapper"
	"../models"
	"../sessions"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "sort"
	_ "strings"
	"time"
)

type result_content struct {
	Page_status  int
	Requirements map[int]*result_req
}

type result_req struct {
	Wording         string
	Solution        string
	ReqType         int
	CompletePercent int
	GoodPercent     int
}

func initRequirements(Requirements map[int]*result_req, currentSession *sessions.Session, db *dbwrapper.DataBaseWrapper) {
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
		FixedSize := len(currentSession.Tests[k].Questions)
		FloatSize := FixedSize
		CompleteCount := 0
		GoodCount := 0
		for _, q := range currentSession.Tests[k].Questions {
			if q.Status == 1 {
				CompleteCount++
				GoodCount++
			}
			if q.Status == 2 {
				CompleteCount++
				GoodCount++
			}
			if q.Status == 3 {
				CompleteCount++
			}
			if q.Status == 4 {
				CompleteCount++
				FloatSize--
			}
		}
		var CompletePercent int
		var GoodPercent int
		if FixedSize != 0 {
			CompletePercent = (int)((float64(CompleteCount) / float64(FixedSize)) * 100.0)
		} else {
			CompletePercent = 100
		}
		if FloatSize != 0 {
			GoodPercent = (int)((float64(GoodCount) / float64(FloatSize)) * 100.0)
		} else {
			GoodPercent = 100
		}
		ReqType := 0
		if CompleteCount == FixedSize {
			if GoodCount == FloatSize {
				ReqType = 2
			} else {
				ReqType = 3
			}
		} else {
			ReqType = 1
		}
		Requirements[k] = &result_req{Req.Wording, Req.Solution, ReqType, CompletePercent, GoodPercent}
	}
}

func GetResultsHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	intr, _ = c.Get("database")
	db := intr.(*dbwrapper.DataBaseWrapper)
	if !currentSession.Fixed {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	rc := result_content{3, make(map[int]*result_req)}
	initRequirements(rc.Requirements, currentSession, db)
	c.HTML(http.StatusOK, "results", rc)
}

type report_content struct {
	CurrentTime     string
	CurrentSession  *sessions.Session
	TypeDescription string
	Requirements    map[int]*result_req
	Clarifications  map[int]*clarify_per_requriement
}

type united_per_clarify struct {
	Wording   string
	Questions map[int]string
}

type clarify_per_requriement struct {
	WithK   map[int]*sessions.StoredQuestion
	WithC   map[int]*united_per_clarify
	IsEmpty bool
	IsK     bool
	IsC     bool
}

func GetReportHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	intr, _ = c.Get("database")
	db := intr.(*dbwrapper.DataBaseWrapper)
	if !currentSession.Fixed {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	rc := report_content{}
	rc.CurrentTime = time.Now().Format("2006-01-02 15:04:05")
	rc.Requirements = make(map[int]*result_req)
	rc.CurrentSession = currentSession
	initRequirements(rc.Requirements, currentSession, db)
	row, err := db.QueryRow("GetTypeListWithId", currentSession.TypeId)
	if err != nil {
		log.Fatal(err)
	}
	err = row.Scan(&rc.TypeDescription)
	if err != nil {
		log.Fatal(err)
	}
	rc.Clarifications = make(map[int]*clarify_per_requriement)
	for key := range rc.Requirements {
		rc.Clarifications[key] = &clarify_per_requriement{}
		rc.Clarifications[key].WithK = make(map[int]*sessions.StoredQuestion)
		rc.Clarifications[key].WithC = make(map[int]*united_per_clarify)
		for k, v := range currentSession.Tests[key].Questions {
			if v.Status == 2 {
				rc.Clarifications[key].WithK[k] = v
			}
			if v.Status == 3 {
				var qrow *sql.Row
				var clarify models.Clarification
				qrow, err = db.QueryRow("GetClarifyWithQue", k)
				if err != nil {
					log.Fatal(err)
				}
				err = qrow.Scan(&clarify.Id, &clarify.Wording)
				if err != nil {
					continue
				}
				_, ok := rc.Clarifications[key].WithC[clarify.Id]
				if !ok {
					rc.Clarifications[key].WithC[clarify.Id] = &united_per_clarify{clarify.Wording, make(map[int]string)}
				}
				rc.Clarifications[key].WithC[clarify.Id].Questions[k] = " " + v.Number
			}
		}
		if len(rc.Clarifications[key].WithK) == 0 && len(rc.Clarifications[key].WithC) == 0 {
			rc.Clarifications[key].IsEmpty = true
		}
		if len(rc.Clarifications[key].WithC) != 0 {
			rc.Clarifications[key].IsC = true
		}
		if len(rc.Clarifications[key].WithK) != 0 {
			rc.Clarifications[key].IsK = true
		}
	}
	c.HTML(http.StatusOK, "report", rc)
}
