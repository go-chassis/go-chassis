package resource

import (
	"errors"
	"net/http"

	rf "github.com/go-chassis/go-chassis/server/restful"
	"math/rand"
	"sync"
)

var num = rand.Intn(100)
var l sync.Mutex

//RestFulMessage is a struct used to implement restful message
type RestFulMessage struct {
}

//Saymessage is used to reply user with his name
func (r *RestFulMessage) DeadLock(b *rf.Context) {
	l.Lock()
	b.Write([]byte("hello world"))
}

//Sayhi is a method used to reply request user with hello world text
func (r *RestFulMessage) Sayhi(b *rf.Context) {
	b.Write([]byte("hello world"))
	return
}

//Sayerror is a method used to reply request user with error
func (r *RestFulMessage) Sayerror(b *rf.Context) {
	b.WriteError(http.StatusInternalServerError, errors.New("test hystric"))
	return
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulMessage) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/lock", ResourceFuncName: "DeadLock"},
		{Method: http.MethodGet, Path: "/sayhimessage", ResourceFuncName: "Sayhi"},
		{Method: http.MethodGet, Path: "/sayerror", ResourceFuncName: "Sayerror"},
	}
}
