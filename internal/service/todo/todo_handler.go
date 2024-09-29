package todo

import (
	"encoding/json"
	"fmt"
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
	var list TodoListPayload
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	if list.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenStr := r.Header.Get("Authorization")[7:]
	email, err := internal.ParseToken(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	user_id, list_id, err := th.storage.CreateListRecord(&list, email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	newList := fmt.Sprintf("{\"New List\": \"list_%v_%v\"}", user_id, list_id)
	w.Write([]byte(newList))
}

func (th *TodoHandler) GetLists(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")[7:]
	email, err := internal.ParseToken(tokenStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	lists, err := th.storage.GetAllLists(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	data, err := json.MarshalIndent(lists, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	w.Write(data)
}

// func (th *TodoHandler) func(w http.ResponseWriter, r *http.Request) {}
