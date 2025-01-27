package routes

import (
	"codeberg.org/sporiff/eigakanban/handlers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes initializes all the routes for the application.
func SetupRoutes(router *gin.Engine, db *pgxpool.Pool) {
	userHandler := handlers.NewUserHandler(db)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			// User routes
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:uuid", userHandler.GetUserByUuid)
			users.PATCH("/:uuid", userHandler.UpdateUser)
			users.POST("/", userHandler.AddUser)
			users.DELETE("/:uuid", userHandler.DeleteUser)
		}
	}
}
