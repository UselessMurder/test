package routes

import (
	"../dbwrapper"
	"../models"
	"../sessions"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type questions_content struct {
	Page_error            int
	Page_status           int
	HasMentor             bool
	HasAddition           bool
	MentorHasAddition     bool
	Question_status       uint8
	Question_compenstaion string
	Quest                 models.Question
	Mentor                models.Question
	Current_session       *sessions.Session
}

func init_question_content(qc *questions_content, currentSession *sessions.Session, db *dbwrapper.DataBaseWrapper) {
	qc.Page_status = 2
	qc.Current_session = currentSession
	var verStr string
	row, err := db.QueryRow("GetQuestionWithId", currentSession.CurrentQuestionId)
	if err != nil {
		log.Fatal(err)
	}
	err = row.Scan(&qc.Quest.Id, &qc.Quest.Number, &qc.Quest.Wording, &qc.Quest.Addition, &verStr)
	if err != nil {
		log.Fatal(err)
	}
	qc.Quest.Verify = strings.Split(verStr, "; ")
	if qc.Quest.Addition != "empty" {
		qc.HasAddition = true
	}
	var srow *sql.Row
	srow, err = db.QueryRow("GetSeniorWithQId", currentSession.CurrentQuestionId)
	if err != nil {
		log.Fatal(err)
	}
	err = srow.Scan(&qc.Mentor.Id, &qc.Mentor.Number, &qc.Mentor.Wording, &qc.Mentor.Addition)
	if err == nil {
		qc.HasMentor = true
		if qc.Mentor.Addition != "empty" {
			qc.MentorHasAddition = true
		}
	}
	qc.Question_status = currentSession.Tests[currentSession.CurrentTestId].Questions[currentSession.CurrentQuestionId].Status
	qc.Question_compenstaion = currentSession.Tests[currentSession.CurrentTestId].Questions[currentSession.CurrentQuestionId].Сompensation
}

func GetQuestionsHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	intr, _ = c.Get("database")
	db := intr.(*dbwrapper.DataBaseWrapper)
	if !currentSession.Fixed {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	currentTestIdStr := c.DefaultQuery("tid", "0")
	currentTestId, _ := strconv.Atoi(currentTestIdStr)
	_, ok := currentSession.Tests[currentTestId]
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/tests")
		return
	}
	currentSession.CurrentTestId = currentTestId
	currentQuestionIdStr := c.DefaultQuery("qid", "0")
	currentQuestionId, _ := strconv.Atoi(currentQuestionIdStr)
	_, ok = currentSession.Tests[currentTestId].Questions[currentQuestionId]
	if !ok {
		if currentSession.Tests[currentTestId].Complete == true && currentSession.Tests[currentTestId].FirstIndex == 0 {
			type ec struct {
				Page_status int
			}
			c.HTML(http.StatusOK, "emptyTest", ec{2})
			return
		}
		currentQuestionId = currentSession.Tests[currentTestId].FirstIndex
	}
	currentSession.CurrentQuestionId = currentQuestionId
	qc := questions_content{}
	init_question_content(&qc, currentSession, db)
	c.HTML(http.StatusOK, "questions", qc)
}

func PostQuestionsHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	intr, _ = c.Get("database")
	db := intr.(*dbwrapper.DataBaseWrapper)
	if !currentSession.Fixed {
		c.Redirect(http.StatusFound, "/")
		return
	}
	if currentSession.CurrentQuestionId == 0 || currentSession.CurrentTestId == 0 {
		c.Redirect(http.StatusFound, "/tests")
		return
	}
	qc := questions_content{}
	statusStr := c.PostForm("QuestionStatus")
	status, err := strconv.Atoi(statusStr)
	log.Println("Status:", status)
	if err != nil || status < 1 || status > 4 {
		init_question_content(&qc, currentSession, db)
		qc.Page_error = 2
		c.HTML(http.StatusOK, "questions", qc)
		return
	}
	if status == 2 {
		KTextStr := c.PostForm("KTextArea")
		currentSession.Tests[currentSession.CurrentTestId].Questions[currentSession.CurrentQuestionId].Сompensation = KTextStr
		if len(strings.Trim(KTextStr, " ")) <= 5 {
			init_question_content(&qc, currentSession, db)
			qc.Page_error = 1
			c.HTML(http.StatusOK, "questions", qc)
			return
		}
	}
	currentSession.Tests[currentSession.CurrentTestId].Questions[currentSession.CurrentQuestionId].Status = uint8(status)
	var keys []int
	keys = make([]int, len(currentSession.Tests[currentSession.CurrentTestId].Questions))
	i := 0
	for k := range currentSession.Tests[currentSession.CurrentTestId].Questions {
		keys[i] = k
		i++
	}
	sort.Ints(keys)
	for k, v := range keys {
		if v == currentSession.CurrentQuestionId {
			if k == (len(keys) - 1) {
				c.Redirect(http.StatusFound, "/tests")
				return
			} else {
				currentSession.CurrentQuestionId = keys[k+1]
				break
			}
		}
	}
	init_question_content(&qc, currentSession, db)
	c.HTML(http.StatusOK, "questions", qc)
}
