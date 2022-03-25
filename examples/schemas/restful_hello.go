package schemas

import (
	"errors"
	"log"
	"net/http"

	"fmt"
	"math/rand"

	rf "github.com/go-chassis/go-chassis/v2/server/restful"
)

var num = rand.Intn(100)

//RestFulHello is a struct used for implementation of restfull hello program
type RestFulHello struct {
}

//Sayhello is a method used to reply user with hello
func (r *RestFulHello) Root(b *rf.Context) {
	b.Write([]byte(fmt.Sprintf("x-forwarded-host %s", b.ReadRequest().Host)))
}

//Sayhello is a method used to reply user with hello
func (r *RestFulHello) Sayhello(b *rf.Context) {
	id := b.ReadPathParameter("userid")
	log.Printf("get user id: " + id)
	log.Printf("get user name: " + b.ReadRequest().Header.Get("user"))
	b.Write([]byte(fmt.Sprintf("user %s from %d", id, num)))
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
	result := struct {
		Name string
	}{}
	err := b.ReadEntity(&result)
	if err != nil {
		b.WriteHeaderAndJSON(http.StatusInternalServerError, result, "application/json")
		return
	}
	result.Name = "hello " + result.Name
	b.WriteJSON(result, "application/json")
	return
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulHello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/", ResourceFunc: r.Root,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodGet, Path: "/sayhello/{userid}", ResourceFunc: r.Sayhello,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayhi", ResourceFunc: r.Sayhi,
			Returns: []*rf.Returns{{Code: 200}}},

		{Method: http.MethodPost, Path: "/sayjson",
			ResourceFunc: r.SayJSON,
			Metadata: map[string]interface{}{
				"tags": []string{"users", "test"},
			},
			Returns: []*rf.Returns{{Code: 200}}},
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

//Sayerror is a method used to reply request user with error
func (r *RestFulMessage) Sayerror(b *rf.Context) {
	b.WriteError(http.StatusInternalServerError, errors.New("test hystric"))
	return
}

//URLPatterns helps to respond for corresponding API calls
func (r *RestFulMessage) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/saymessage/{name}", ResourceFunc: r.Saymessage},
		{Method: http.MethodPost, Path: "/sayhimessage", ResourceFunc: r.Sayhi},
		{Method: http.MethodGet, Path: "/sayerror", ResourceFunc: r.Sayhi},
	}
}

//Hello is a struct used for implementation of restfull hello program
type Hello struct{}

//Hello
func (r *Hello) Hello(b *rf.Context) { b.Write([]byte("hi from hello")) }

//URLPatterns helps to respond for corresponding API calls
func (r *Hello) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/hello", ResourceFunc: r.Hello},
	}
}

//Legacy is a struct
type Legacy struct{}

//Do
func (r *Legacy) Do(b *rf.Context) { b.Write([]byte("hello from legacy")) }

//URLPatterns helps to respond for corresponding API calls
func (r *Legacy) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/legacy", ResourceFunc: r.Do},
	}
}

//Legacy is a struct
type Admin struct{}

//Do
func (r *Admin) Do(b *rf.Context) { b.Write([]byte("hello from admin")) }

//URLPatterns helps to respond for corresponding API calls
func (r *Admin) URLPatterns() []rf.Route {
	return []rf.Route{
		{Method: http.MethodGet, Path: "/admin", ResourceFunc: r.Do},
	}
}
