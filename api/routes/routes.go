package routes

import (
	"codeberg.org/sporiff/eigakanban/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes initializes all the routes for the application.
func SetupRoutes(router *gin.Engine, db *pgxpool.Pool) {
	userHandler := handlers.NewUserHandler(db)

	// User routes
	router.GET("/users", userHandler.GetAllUsers)
	router.GET("/users/:uuid", userHandler.GetUserByUuid)
	router.PATCH("/users/:uuid", userHandler.UpdateUser)
	router.POST("/users", userHandler.AddUser)
	router.DELETE("/users/:uuid", userHandler.DeleteUser)
}
