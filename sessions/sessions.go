package sessions

import (
	"../hashgenerator"
	"errors"
	"net/http"
	"time"
)

const (
	COOKIE_NAME = "testsessionid"
)

type Session struct {
	SessionHash    string
	IsTestRunning  bool
	TestType       uint8
	QuestionNumber uint16
	ExpireTime     time.Time
}

type SessionManager struct {
	sessions   map[string]*Session
	setChan    chan *Session
	getChan    chan string
	tubeChan   chan *Session
	removeChan chan string
	expireChan chan struct{}
	doneChan   chan struct{}
}

func (sm *SessionManager) OpenSessionManager() {
	sm.sessions = make(map[string]*Session)
	sm.setChan = make(chan *Session)
	sm.tubeChan = make(chan *Session)
	sm.getChan = make(chan string)
	sm.removeChan = make(chan string)
	sm.doneChan = make(chan struct{})
	sm.expireChan = make(chan struct{})

	go func() {
		for {
			select {
			case currentSession := <-sm.setChan:
				sm.sessions[currentSession.SessionHash] = currentSession
			case sessionId := <-sm.getChan:
				sm.tubeChan <- sm.sessions[sessionId]
			case sessionId := <-sm.removeChan:
				delete(sm.sessions, sessionId)
			case <-sm.expireChan:
				currentTime := time.Now()
				for _, currentSession := range sm.sessions {
					if currentTime.After(currentSession.ExpireTime) {
						delete(sm.sessions, currentSession.SessionHash)
					}
				}
			case <-sm.doneChan:
				return
			}
		}
	}()

	go func() {
		for {
			time.Sleep(5 * time.Minute)
			sm.expireChan <- struct{}{}
		}
	}()
}

func (sm *SessionManager) CloseSessionManager() {
	sm.doneChan <- struct{}{}
}

func (sm *SessionManager) SetSession(currentSession *Session) {
	sm.setChan <- currentSession
}

func (sm *SessionManager) GetSession(sessionsHash string) (error, *Session) {
	sm.getChan <- sessionsHash
	currentSession := <-sm.tubeChan
	var err error
	if currentSession == nil {
		err = errors.New("Invalid session!")
	} else {
		currentSession.ExpireTime = time.Now().Add(24 * time.Hour)
	}

	return err, currentSession
}

func (sm *SessionManager) GetCookie(r *http.Request, w http.ResponseWriter) string {

	cookie, err := r.Cookie(COOKIE_NAME)

	if err != nil {

		hash, _ := hashgenerator.GenerateHash28(time.Now().String(), "User")
		t := time.Now().Add(24 * time.Hour)

		sm.setChan <- &Session{hash, false, 0, 0, t}

		cookie = &http.Cookie{
			Name:    COOKIE_NAME,
			Value:   hash,
			Expires: t,
		}

	} else {
		cookie.Expires = time.Now().Add(24 * time.Hour)
	}

	http.SetCookie(w, cookie)

	return cookie.Value
}
