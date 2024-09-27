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
	subrouter := router.PathPrefix("/api/v1/").Subrouter()

	RegisterUserRoutes(subrouter, s.db)
	RegisterTodoRoutes(subrouter, s.db)

	router.Use(LogMW, AuthMW)
	log.Printf("Server is up and running on port: %v", s.addr[1:])
	return http.ListenAndServe(s.addr, router)
}
