package server

import (
	"database/sql"

	"github.com/AHMED-D007A/Todo-List-API/internal/service/todo"
	"github.com/AHMED-D007A/Todo-List-API/internal/service/user"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, db *sql.DB) {
	userHandler := user.NewUserHanlder(user.NewUserStorage(db))

	router.HandleFunc("/register", userHandler.RegisterNewUserHandler).Methods("POST")
	router.HandleFunc("/login", userHandler.VerifiyUserHandler).Methods("POST")
}

func RegisterTodoRoutes(router *mux.Router, db *sql.DB) {
	todoHandler := todo.NewTodoHandler(todo.NewTodoStorage(db))

	/* a route where a user create list with these fields:(title) */
	router.HandleFunc("/lists", todoHandler.CreateNewList).Methods("POST")
	/* a route where the client gets all lists that the user has created before. */
	router.HandleFunc("/lists", todoHandler.GetLists).Methods("GET")

	/* in every list that the user has created he can get all the todo_items inside the list */
	router.HandleFunc("/lists/{signature}/todos", todoHandler.GetItems).Methods("GET")
	/* in every list that the user has created he can created a new task(todo_item) */
	router.HandleFunc("/lists/{signature}/todos", todoHandler.CreateNewItem).Methods("POST")
	/* in every list that the user has created he can updated the status
	just make sure to send the title, description status every time you are going to update */
	router.HandleFunc("/lists/{signature}/todos/{todo_id}", todoHandler.UpdateItem).Methods("PUT")
	/* in every list that the user has created he can delete a todo_item that has been created before */
	router.HandleFunc("/lists/{signature}/todos/{todo_id}", todoHandler.DeleteItem).Methods("DELETE")
}
