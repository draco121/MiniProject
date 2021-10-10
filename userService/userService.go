package userService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type User struct {
	Id      int    `json:"Id"`
	Name    string `json:"Name"`
	Email   string `json:"Email"`
	Age     int    `json:"Age"`
	Contact string `json:"Contact"`
}

type userHandler struct {
	details map[string]User
}

func NewUserHandler() *userHandler {
	return &userHandler{
		details: map[string]User{
			// }
			"1": User{
				Id:      1,
				Name:    "Srijan",
				Email:   "srijan@gmail.com",
				Age:     23,
				Contact: "8223821911",
			},
		},
	}

}

func (u *userHandler) Default(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		u.get(rw, r)
	case "POST":
		u.post(rw, r)
	case "PUT":
		u.put(rw, r)
	case "DELETE":
		u.delete(rw, r)
	}
}

func (u *userHandler) get(rw http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")
	user, ok := u.details[parts[2]]

	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	rw.Header().Add("content-type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(data)
}

func (u *userHandler) post(rw http.ResponseWriter, r *http.Request) {

	if r.Body == http.NoBody {
		rw.WriteHeader(http.StatusBadRequest)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	id := len(u.details)
	user.Id = id + 1

	u.details[fmt.Sprint(user.Id)] = user

	rw.WriteHeader(http.StatusCreated)
}

func (u *userHandler) put(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var user User
	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	u.details[fmt.Sprint(user.Id)] = user

	rw.WriteHeader(http.StatusOK)
}

func (u *userHandler) delete(rw http.ResponseWriter, r *http.Request) {

	parts := strings.Split(r.URL.String(), "/")
	delete(u.details, parts[2])

	rw.WriteHeader(http.StatusOK)
}
