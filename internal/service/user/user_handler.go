package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AHMED-D007A/Todo-List-API/internal"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	storage *UserStorage
}

func NewUserHanlder(storage *UserStorage) *UserHandler {
	return &UserHandler{
		storage: storage,
	}
}

func (uh *UserHandler) RegisterNewUserHandler(w http.ResponseWriter, r *http.Request) {
	var userpayload UserPayload
	err := json.NewDecoder(r.Body).Decode(&userpayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	userpayload.Password, err = bcrypt.GenerateFromPassword([]byte(userpayload.Password), 14)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := internal.CreateToken(userpayload.Name, userpayload.Email)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = uh.storage.RegisterNewUserStorage(&userpayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	response := fmt.Sprintf("{\n\t\"token\": \"%v\"\n}", token)

	w.Write([]byte(response))
}

func (uh *UserHandler) VerifiyUserHandler(w http.ResponseWriter, r *http.Request) {
	var userpayload UserPayload
	err := json.NewDecoder(r.Body).Decode(&userpayload)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := uh.storage.VerifiyUserStorage(&userpayload)
	if err != nil {
		log.Print(err.Error())
		if err.Error() == "not Found" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("This Email is Not Registered"))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(user.Password, userpayload.Password)
	if err != nil {
		w.Write([]byte("The email or the password are wrong."))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := internal.CreateToken(userpayload.Name, userpayload.Email)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("{\"token\": \"%v\"}", token)

	w.Write([]byte(response))
}
