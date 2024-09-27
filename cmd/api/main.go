package main

import (
	"fmt"
	"log"

	"github.com/AHMED-D007A/Todo-List-API/internal"
	"github.com/AHMED-D007A/Todo-List-API/internal/database"
	"github.com/AHMED-D007A/Todo-List-API/internal/server"
)

func main() {
	connStr := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		internal.Envs.DB_HOST, internal.Envs.DB_PORT, internal.Envs.DB_USERNAME,
		internal.Envs.DB_PASSWORD, internal.Envs.DB_NAME)
	db := database.NewDBConnection(connStr)
	defer db.Close()

	server := server.NewAPIServer(":4000", db)
	if err := server.Run(); err != nil {
		log.Panic(err.Error())
	}
}
