package routes

import (
	"codeberg.org/sporiff/eigakanban/handlers"
	"codeberg.org/sporiff/eigakanban/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes initializes all the routes for the application.
func SetupRoutes(router *gin.Engine, db *pgxpool.Pool) {
	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	boardsHandler := handlers.NewBoardsHandler(db)

	authMiddlewareHandler := middleware.NewAuthMiddlewareHandler(db)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		users.Use(authMiddlewareHandler.AuthRequired())
		{
			// User routes
			users.GET("/", userHandler.GetAllUsers)
			users.GET("/:uuid", userHandler.GetUserByUuid)
			users.PATCH("/:uuid", userHandler.UpdateUser)
			users.DELETE("/:uuid", userHandler.DeleteUser)
			users.GET("/:uuid/boards", boardsHandler.GetBoardsForUser)
		}
		boards := v1.Group("/boards")
		boards.Use(authMiddlewareHandler.AuthRequired())
		{
			// Board routes
			boards.GET("/", boardsHandler.GetAllBoards)
			boards.GET("/:uuid", boardsHandler.GetBoardByUuid)
			boards.POST("/", boardsHandler.AddBoard)
			boards.PATCH("/:uuid", boardsHandler.UpdateBoard)
			boards.DELETE("/:uuid", boardsHandler.DeleteBoard)
		}
		auth := v1.Group("/auth")
		{
			// Auth routes
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/login", authHandler.LoginUser)
			auth.POST("/logout", authHandler.LogoutUser)
		}
	}
}
