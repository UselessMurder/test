package routes

import (
	"../dbwrapper"
	"../models"
	"../sessions"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type index_content struct {
	Page_error      int
	Page_status     int
	Current_session *sessions.Session
	List_types      map[int]models.AuditType
}

func indexHandler(c *gin.Context, ic *index_content) {
	if ic.Current_session.Fixed {
		ic.Page_status = 1
	}
	c.HTML(http.StatusOK, "index", ic)
}

func GetIndexHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	intr, _ = c.Get("database")
	db := intr.(*dbwrapper.DataBaseWrapper)
	rows, err := db.Query("GetTypeList")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	ic := index_content{0, 0, currentSession, make(map[int]models.AuditType)}
	for rows.Next() {
		var currentType models.AuditType
		err := rows.Scan(&currentType.Id, &currentType.Name, &currentType.Description)
		if err != nil {
			log.Fatal(err)
		}
		ic.List_types[currentType.Id] = currentType
	}
	indexHandler(c, &ic)
	return
}

func PostIndexHandler(c *gin.Context) {
	intr, _ := c.Get("currentSession")
	currentSession := intr.(*sessions.Session)
	currentSession.InOrganizationName = c.PostForm("InOrganizationName")
	currentSession.CommercialDesignation = c.PostForm("CommercialDesignation")
	currentSession.InContactPerson = c.PostForm("InContactPerson")
	currentSession.InContactPersonPost = c.PostForm("InContactPersonPost")
	currentSession.InPhone = c.PostForm("InPhone")
	currentSession.InEmail = c.PostForm("InEmail")
	currentSession.InAddress = c.PostForm("InAddress")
	currentSession.InCity = c.PostForm("InCity")
	currentSession.InState = c.PostForm("InState")
	currentSession.InCountry = c.PostForm("InCountry")
	currentSession.InIndex = c.PostForm("InIndex")
	currentSession.InURL = c.PostForm("InURL")
	currentSession.OutOrganizationName = c.PostForm("OutOrganizationName")
	currentSession.OutContactPerson = c.PostForm("OutContactPerson")
	currentSession.OutContactPersonPost = c.PostForm("OutContactPersonPost")
	currentSession.OutPhone = c.PostForm("OutPhone")
	currentSession.OutEmail = c.PostForm("OutEmail")
	currentSession.OutAddress = c.PostForm("OutAddress")
	currentSession.OutCity = c.PostForm("OutCity")
	currentSession.OutState = c.PostForm("OutState")
	currentSession.OutCountry = c.PostForm("OutCountry")
	currentSession.OutIndex = c.PostForm("OutIndex")
	currentSession.OutURL = c.PostForm("OutURL")
	errCode := validate(currentSession)
	types := make(map[int]models.AuditType)
	intr, _ = c.Get("database")
	db := intr.(*dbwrapper.DataBaseWrapper)
	rows, err := db.Query("GetTypeList")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var currentType models.AuditType
		err := rows.Scan(&currentType.Id, &currentType.Name, &currentType.Description)
		if err != nil {
			log.Fatal(err)
		}
		types[currentType.Id] = currentType
	}
	var TypeId int
	TypeId, err = strconv.Atoi(c.PostForm("TypeId"))
	if TypeId < 1 || TypeId > len(types) {
		errCode = 24
	} else {
		currentSession.TypeId = TypeId
	}
	if errCode != 0 {
		indexHandler(c, &index_content{errCode, 0, currentSession, types})
	} else {
		currentSession.Fixed = true
		for k := range currentSession.Tests {
			delete(currentSession.Tests, k)
		}
		currentSession.CurrentTestId = 0
		currentSession.CurrentQuestionId = 0
		var rrows *sql.Rows
		rrows, err = db.Query("GetRequrinmentsWithId", currentSession.TypeId)
		if err != nil {
			log.Fatal(err)
		}
		defer rrows.Close()
		for rrows.Next() {
			var reqId int
			err := rrows.Scan(&reqId)
			if err != nil {
				log.Fatal(err)
			}
			currentSession.Tests[reqId] = &sessions.StoredTest{false, 0, make(map[int]*sessions.StoredQuestion)}
		}
		for k := range currentSession.Tests {
			var qrows *sql.Rows
			qrows, err := db.Query("GetQuestionsWithReq", k)
			if err != nil {
				log.Fatal(err)
			}
			defer qrows.Close()
			lock := false
			for qrows.Next() {
				var queId int
				var queNumber string
				err := qrows.Scan(&queId, &queNumber)
				if err != nil {
					log.Fatal(err)
				}
				if !lock {
					currentSession.Tests[k].FirstIndex = queId
					lock = true
				}
				currentSession.Tests[k].Questions[queId] = &sessions.StoredQuestion{0, "", queNumber}
			}
		}
		for k := range currentSession.Tests {
			if currentSession.Tests[k].FirstIndex == 0 {
				currentSession.Tests[k].Complete = true
			}
		}
		c.Redirect(http.StatusFound, "/tests")
	}
}

func validate(currentSession *sessions.Session) int {
	ok, _ := regexp.MatchString(`[a-zA-Z0-9а-яёА-ЯЁ .,!?&()*$#@+=-_><'"]{1,30}$`, currentSession.InOrganizationName)
	if !ok {
		return 1
	}
	ok, _ = regexp.MatchString(`[a-zA-Z0-9а-яёА-ЯЁ .,!?&()*$#@+=-_><'"]{1,30}$`, currentSession.CommercialDesignation)
	if !ok {
		return 2
	}
	ok, _ = regexp.MatchString(`^[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+$`, currentSession.InContactPerson)
	if !ok {
		return 3
	}
	ok, _ = regexp.MatchString(`[А-ЯЁa-яё 1-9]{1,20}`, currentSession.InContactPersonPost)
	if !ok {
		return 4
	}
	ok, _ = regexp.MatchString(`^\+?\d{1,3}?[- .]?\(?(?:\d{2,3})\)?[- .]?\d\d\d[- .]?\d\d\d\d$`, currentSession.InPhone)
	if !ok {
		return 5
	}
	ok, _ = regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, currentSession.InEmail)
	if !ok {
		return 6
	}
	ok, _ = regexp.MatchString(`^(.+)\s+(\S+?)(-(\d+))?$`, currentSession.InAddress)
	if !ok {
		return 7
	}
	ok, _ = regexp.MatchString(`[a-zA-Zа-яёА-ЯЁ .]{1,30}`, currentSession.InCity)
	if !ok {
		return 8
	}
	ok, _ = regexp.MatchString(`[a-zA-Zа-яёА-ЯЁ .]{1,30}`, currentSession.InState)
	if !ok {
		return 9
	}
	ok, _ = regexp.MatchString(`[a-zA-Zа-яёА-ЯЁ .]{1,30}`, currentSession.InCountry)
	if !ok {
		return 10
	}
	ok, _ = regexp.MatchString(`^\d{6}$`, currentSession.InIndex)
	if !ok {
		return 11
	}
	ok, _ = regexp.MatchString(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`, currentSession.InURL)
	if !ok {
		return 12
	}
	ok, _ = regexp.MatchString(`[a-zA-Z0-9а-яёА-ЯЁ .,!?&()*$#@+=-_><'"]{1,30}$`, currentSession.OutOrganizationName)
	if !ok {
		return 13
	}
	ok, _ = regexp.MatchString(`^[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+$`, currentSession.OutContactPerson)
	if !ok {
		return 14
	}
	ok, _ = regexp.MatchString(`[А-ЯЁa-яё 1-9]{1,20}`, currentSession.OutContactPersonPost)
	if !ok {
		return 15
	}
	ok, _ = regexp.MatchString(`^\+?\d{1,3}?[- .]?\(?(?:\d{2,3})\)?[- .]?\d\d\d[- .]?\d\d\d\d$`, currentSession.OutPhone)
	if !ok {
		return 16
	}
	ok, _ = regexp.MatchString(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`, currentSession.OutEmail)
	if !ok {
		return 17
	}
	ok, _ = regexp.MatchString(`^(.+)\s+(\S+?)(-(\d+))?$`, currentSession.OutAddress)
	if !ok {
		return 18
	}
	ok, _ = regexp.MatchString(`[a-zA-Zа-яёА-ЯЁ .]{1,30}`, currentSession.OutCity)
	if !ok {
		return 19
	}
	ok, _ = regexp.MatchString(`[a-zA-Zа-яёА-ЯЁ .]{1,30}`, currentSession.OutState)
	if !ok {
		return 20
	}
	ok, _ = regexp.MatchString(`[a-zA-Zа-яёА-ЯЁ .]{1,30}`, currentSession.OutCountry)
	if !ok {
		return 21
	}
	ok, _ = regexp.MatchString(`^\d{6}$`, currentSession.OutIndex)
	if !ok {
		return 22
	}
	ok, _ = regexp.MatchString(`^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`, currentSession.OutURL)
	if !ok {
		return 23
	}
	return 0
}
