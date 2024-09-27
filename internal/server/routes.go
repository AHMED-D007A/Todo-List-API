package server

import (
	"database/sql"

	"github.com/AHMED-D007A/Todo-List-API/internal/service/user"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router, db *sql.DB) {
	userHandler := user.NewUserHanlder(user.NewUserStorage(db))

	router.HandleFunc("/register", userHandler.RegisterNewUserHandler).Methods("POST")
	router.HandleFunc("/login", userHandler.VerifiyUserHandler).Methods("POST")
}

func RegisterTodoRoutes(router *mux.Router, db *sql.DB) {}
