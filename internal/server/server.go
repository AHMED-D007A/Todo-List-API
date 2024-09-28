package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	User_subrouter := router.PathPrefix("/api/v1/").Subrouter()
	Todo_subrouter := router.PathPrefix("/api/v1/").Subrouter()
	router.Use(LogMW)
	Todo_subrouter.Use(AuthMW)

	RegisterUserRoutes(User_subrouter, s.db)
	RegisterTodoRoutes(Todo_subrouter, s.db)

	log.Printf("Server is up and running on port: %v", s.addr[1:])
	return http.ListenAndServe(s.addr, router)
}
