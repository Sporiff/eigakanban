package routes

import (
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/handlers"
	"codeberg.org/sporiff/eigakanban/middleware"
	"codeberg.org/sporiff/eigakanban/services"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SetupRoutes initializes all the routes for the application.
func SetupRoutes(router *gin.Engine, db *pgxpool.Pool) {
	q := queries.New(db)

	authService := services.NewAuthService(q)
	usersService := services.NewUsersService(q)
	statusesService := services.NewStatusesService(q)
	itemsService := services.NewItemsService(q)
	listItemsService := services.NewListItemsService(q)

	authHandler := handlers.NewAuthHandler(authService)
	usersHandler := handlers.NewUsersHandler(usersService)
	statusesHandler := handlers.NewStatusesHandler(statusesService)
	itemsHandler := handlers.NewItemsHandler(itemsService)
	listItemsHandler := handlers.NewListItemsHandler(listItemsService)

	authMiddlewareHandler := middleware.NewAuthMiddlewareHandler(db)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		users.Use(authMiddlewareHandler.AuthRequired())
		{
			// User routes
			users.GET("/", usersHandler.GetAllUsers)
			users.GET("/:uuid", usersHandler.GetUserByUuid)
			users.PATCH("/:uuid", usersHandler.UpdateUser)
			users.DELETE("/:uuid", usersHandler.DeleteUser)
		}
		items := v1.Group("/items")
		items.Use(authMiddlewareHandler.AuthRequired())
		{
			// Items routes
			items.GET("/", itemsHandler.GetAllItems)
			items.GET("/:uuid", itemsHandler.GetItemByUuid)
			items.POST("/", itemsHandler.AddItem)
			items.PATCH("/:uuid", itemsHandler.UpdateItem)
		}
		lists := v1.Group("/lists")
		lists.Use(authMiddlewareHandler.AuthRequired())
		{
			// Lists handlers
			lists.GET("/:uuid", listItemsHandler.GetListItemsForList)
		}
		listItems := v1.Group("/list_items")
		listItems.Use(authMiddlewareHandler.AuthRequired())
		{
			// List items handlers
			listItems.GET("/", listItemsHandler.GetAllListItems)
		}
		statuses := v1.Group("/statuses")
		statuses.Use(authMiddlewareHandler.AuthRequired())
		{
			// Statuses handlers
			statuses.POST("/", statusesHandler.AddStatus)
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
