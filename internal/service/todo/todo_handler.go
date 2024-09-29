package todo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AHMED-D007A/Todo-List-API/internal"
	"github.com/gorilla/mux"
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
	w.WriteHeader(http.StatusCreated)
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

func (th *TodoHandler) CreateNewItem(w http.ResponseWriter, r *http.Request) {
	var item TodoItemPayload
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	if item.Title == "" || item.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)

	todoITem, err := th.storage.CreateItemRecord(vars["signature"], &item)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	data, err := json.MarshalIndent(todoITem, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (th *TodoHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var payload TodoItemPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	if payload.Title == "" && payload.Description == "" && payload.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	item, err := th.storage.UpdateItemRecord(vars["signature"], vars["todo_id"], payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	data, err := json.MarshalIndent(item, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(data)
}

func (th *TodoHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := th.storage.DeleteItemRecord(vars["signature"], vars["todo_id"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (th *TodoHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	items, err := th.storage.GetAllItems(vars["signature"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	data, err := json.MarshalIndent(items, "", "\t")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err.Error())
		return
	}

	w.Write(data)
}

// func (th *TodoHandler) func(w http.ResponseWriter, r *http.Request) {}
