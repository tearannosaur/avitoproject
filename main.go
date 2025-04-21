package main

import (
	"log"

	"github.com/project/appserver"
	"github.com/project/database"
	"github.com/project/pkg/handler"
)

func main() {

	db := database.DbInit()
	h := handler.NewHandler(db)
	router := h.InitRoutes()
	srv := new(appserver.Server)
	if err := srv.Run("8080", router); err != nil {
		log.Fatalf("error for running server:%s", err.Error())
	}

}
