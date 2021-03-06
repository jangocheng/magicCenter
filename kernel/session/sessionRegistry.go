package session

import (
	"log"
	"net/http"
	"time"

	"muidea.com/magicCenter/common"
	common_const "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/foundation/util"
)

var sessionCookieID = "session_id"

func init() {
	sessionCookieID = createUUID()
}

func createUUID() string {
	return util.RandomAlphanumeric(32)
}

type sessionRegistryImpl struct {
	commandChan commandChanImpl
}

// CreateSessionRegistry 创建Session仓库
func CreateSessionRegistry() common.SessionRegistry {
	impl := sessionRegistryImpl{}
	impl.commandChan = make(commandChanImpl)
	go impl.commandChan.run()
	go impl.checkTimer()

	return &impl
}

// GetSession 获取Session对象
func (sm *sessionRegistryImpl) GetSession(w http.ResponseWriter, r *http.Request) common.Session {
	var userSession common.Session

	sessionID := ""
	cookie, err := r.Cookie(sessionCookieID)
	if err == nil {
		log.Printf("url sessionID:%s, cookie sessionID:%s, sessionCookieID:%s", sessionID, cookie.Value, common_const.SessionID)
		sessionID = cookie.Value
	}
	urlSession := r.URL.Query().Get(common_const.SessionID)
	if len(urlSession) > 0 {
		sessionID = urlSession
	}

	cur, found := sm.FindSession(sessionID)
	if !found {
		log.Printf("can\\'t find session,create new session, sessionID:%s", sessionID)
		sessionID := createUUID()
		userSession = sm.CreateSession(sessionID)
	} else {
		log.Printf("find exist session, sessionID:%s", sessionID)
		userSession = cur
	}

	// 存入cookie,使用cookie存储
	sessionCookie := http.Cookie{Name: sessionCookieID, Value: userSession.ID(), Path: "/"}
	http.SetCookie(w, &sessionCookie)

	r.AddCookie(&sessionCookie)

	return userSession
}

// CreateSession 新建Session
func (sm *sessionRegistryImpl) CreateSession(sessionID string) common.Session {
	session := sessionImpl{id: sessionID, context: make(map[string]interface{}), registry: sm}

	session.refresh()

	sm.commandChan.insert(session)

	return &session
}

func (sm *sessionRegistryImpl) FindSession(sessionID string) (common.Session, bool) {
	session, found := sm.commandChan.find(sessionID)
	return &session, found
}

// UpdateSession 更新Session
func (sm *sessionRegistryImpl) UpdateSession(session common.Session) bool {
	cur, found := sm.commandChan.find(session.ID())
	if !found {
		return false
	}

	for _, key := range session.OptionKey() {
		cur.context[key], _ = session.GetOption(key)
	}

	return sm.commandChan.update(cur)
}

func (sm *sessionRegistryImpl) FlushSession(session common.Session) {
	sm.commandChan.flush(session.ID())
}

func (sm *sessionRegistryImpl) checkTimer() {
	timeOutTimer := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-timeOutTimer.C:
			sm.commandChan.checkTimeOut()
		}
	}
}

func (sm *sessionRegistryImpl) insert(session sessionImpl) {
	sm.commandChan.insert(session)
}

func (sm *sessionRegistryImpl) delete(id string) {
	sm.commandChan.remove(id)
}

func (sm *sessionRegistryImpl) find(id string) (sessionImpl, bool) {
	return sm.commandChan.find(id)
}

func (sm *sessionRegistryImpl) count() int {
	return sm.commandChan.count()
}

func (sm *sessionRegistryImpl) update(session sessionImpl) bool {
	return sm.commandChan.update(session)
}

type commandData struct {
	action commandAction
	value  interface{}
	result chan<- interface{}
	data   chan<- map[string]interface{}
}

type commandAction int

const (
	insert commandAction = iota
	remove
	update
	find
	checkTimeOut
	length
	flush
	end
)

type findResult struct {
	value interface{}
	found bool
}

type commandChanImpl chan commandData

func (right commandChanImpl) insert(session sessionImpl) {
	log.Printf("insert session, id:%s", session.id)

	right <- commandData{action: insert, value: session}
}

func (right commandChanImpl) remove(id string) {
	log.Printf("delete session, id:%s", id)
	right <- commandData{action: remove, value: id}
}

func (right commandChanImpl) update(session sessionImpl) bool {
	log.Printf("update session, id:%s", session.id)

	reply := make(chan interface{})
	right <- commandData{action: update, value: session, result: reply}

	result := (<-reply).(bool)
	return result
}

func (right commandChanImpl) find(id string) (sessionImpl, bool) {
	log.Printf("find session by id, id:%s", id)
	reply := make(chan interface{})
	right <- commandData{action: find, value: id, result: reply}

	result := (<-reply).(findResult)

	if result.found {
		return result.value.(sessionImpl), result.found
	}

	return sessionImpl{}, false
}

func (right commandChanImpl) count() int {
	log.Print("count session")

	reply := make(chan interface{})
	right <- commandData{action: length, result: reply}

	result := (<-reply).(int)
	return result
}

func (right commandChanImpl) flush(id string) {
	reply := make(chan interface{})
	right <- commandData{action: flush, value: id, result: reply}

	<-reply
}

func (right commandChanImpl) run() {
	sessionContextMap := make(map[string]interface{})
	for command := range right {
		switch command.action {
		case insert:
			session := command.value.(sessionImpl)
			sessionContextMap[session.id] = &session
		case remove:
			id := command.value.(string)
			delete(sessionContextMap, id)
		case update:
			session := command.value.(sessionImpl)
			_, found := sessionContextMap[session.id]
			if found {
				sessionContextMap[session.id] = &session
			}

			command.result <- found
		case find:
			id := command.value.(string)
			session := sessionImpl{}
			cur, found := sessionContextMap[id]
			if found {
				cur.(*sessionImpl).refresh()
				session = *(cur.(*sessionImpl))
			}
			command.result <- findResult{session, found}
		case checkTimeOut:
			removeList := []string{}
			for k, v := range sessionContextMap {
				session := v.(*sessionImpl)
				if session.timeOut() {
					removeList = append(removeList, k)
				}
			}

			for key := range removeList {
				delete(sessionContextMap, removeList[key])
			}
		case length:
			command.result <- len(sessionContextMap)
		case flush:
			command.result <- true
		case end:
			close(right)
			command.data <- sessionContextMap
		}
	}

	log.Print("session manager sessionImpl exit")
}

func (right commandChanImpl) close() map[string]interface{} {
	reply := make(chan map[string]interface{})
	right <- commandData{action: end, data: reply}
	return <-reply
}

func (right commandChanImpl) checkTimeOut() {
	right <- commandData{action: checkTimeOut}
}
