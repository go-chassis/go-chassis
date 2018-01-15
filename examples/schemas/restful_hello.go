package schemas

import (
	"errors"
	"log"
	"net/http"

	rf "github.com/ServiceComb/go-chassis/server/restful"
)

//RestFulHello is a struct used for implementation of restfull hello program
type RestFulHello struct {
}

//Sayhello is a method used to reply user with hello
func (r *RestFulHello) Sayhello(b *rf.Context) {
	id := b.ReadPathParameter("userid")
	log.Printf("get user id: " + id)
	b.Write([]byte("get user id: " + id))
}

//Sayhi is a method used to reply user with hello world text
func (r *RestFulHello) Sayhi(b *rf.Context) {
	result := struct {
		Name string
	}{}
	err := b.ReadEntity(&result)
	if err != nil {
		b.Write([]byte(err.Error() + ":hello world"))
		return
	}
	b.Write([]byte(result.Name + ":hello world"))
	return
}

// SayJSON is a method used to reply user hello in json format
func (r *RestFulHello) SayJSON(b *rf.Context) {
	reslut := struct {
		Name string
	}{}
	err := b.ReadEntity(&reslut)
	if err != nil {
		b.WriteHeaderAndJSON(http.StatusInternalServerError, reslut, "application/json")
		return
	}
	reslut.Name = "hello " + reslut.Name
	b.WriteJSON(reslut, "application/json")
	return
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{http.MethodGet, "/sayhello/{userid}", "Sayhello"},
		{http.MethodPost, "/sayhi", "Sayhi"},
		{http.MethodPost, "/sayjson", "SayJSON"},
	}
}

//RestFulMessage is a struct used to implement restful message
type RestFulMessage struct {
}

//Saymessage is used to reply user with his name
func (r *RestFulMessage) Saymessage(b *rf.Context) {
	id := b.ReadPathParameter("name")

	b.Write([]byte("get name: " + id))
}

//Sayhi is a method used to reply request user with hello world text
func (r *RestFulMessage) Sayhi(b *rf.Context) {
	reslut := struct {
		Name string
	}{}
	err := b.ReadEntity(&reslut)
	if err != nil {
		b.Write([]byte(err.Error() + ":hello world"))
		return
	}
	b.Write([]byte(reslut.Name + ":hello world"))
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
		{http.MethodGet, "/saymessage/{name}", "Saymessage"},
		{http.MethodPost, "/sayhimessage", "Sayhi"},
		{http.MethodGet, "/sayerror", "Sayerror"},
	}
}
