package main

import (
	"codeberg.org/sporiff/eigakanban/config"
	"codeberg.org/sporiff/eigakanban/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	dbConfig := config.LoadDBConfig()

	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatal("Couldn't connect to the database: %v", err)
	}
	defer db.Close()

	router := gin.Default()
	routes.SetupRoutes(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Couldn't start server: %v", err)
	}
}
