package todo

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AHMED-D007A/Todo-List-API/internal"
)

type TodoHandler struct {
	storage TodoStorage
}

func NewTodoHandler(storage *TodoStorage) *TodoHandler {
	return &TodoHandler{
		storage: *storage,
	}
}

func (th *TodoHandler) CreateNewList(w http.ResponseWriter, r *http.Request) {
	var list TodoList
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	tokenStr := r.Header.Get("Authorization")[7:]
	email, err := internal.ParseToken(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	err = th.storage.CreateListRecord(&list, email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}
}

// func (th *TodoHandler) func(w http.ResponseWriter, r *http.Request) {}
