package main

import (
	"codeberg.org/sporiff/eigakanban/config"
	_ "codeberg.org/sporiff/eigakanban/docs"
	"codeberg.org/sporiff/eigakanban/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

//	@title			eigakanban API
//	@version		1.0
//	@description	The REST API for the eigakanban server
//	@BasePath		/api/v1

//	@contact.name	eigakanban developers
//	@contact.url	https://codeberg.org/sporiff/eigakanban/issues
//	@contact.email	~sporiff/eigakanban@lists.sr.ht

//	@license.name	AGPL3 or Later
//	@license.url	https://codeberg.org/sporiff/eigakanban/src/branch/main/LICENSE

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {

	dbConfig := config.LoadDBConfig()

	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %v", err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(cors.Default())
	routes.SetupRoutes(router, db, tmdbClient)

	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Couldn't start server: %v", err)
	}
}
