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
	superUserMiddlewareHandler := middleware.NewSuperUserMiddlewareHandler(db)

	v1 := router.Group("/api/v1")
	{
		// Unauthenticated routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/login", authHandler.LoginUser)
			auth.POST("/logout", authHandler.LogoutUser)
		}

		items := v1.Group("/items")
		{
			items.GET("/", itemsHandler.GetAllItems)
			items.GET("/:uuid", itemsHandler.GetItemByUuid)
		}

		lists := v1.Group("/lists")
		{
			lists.GET("/:uuid", listItemsHandler.GetListItemsForList)
		}

		listItems := v1.Group("/lists/:uuid")
		{
			listItems.GET("/", listItemsHandler.GetAllListItems)
		}

		// Authenticated routes
		users := v1.Group("/users/:uuid")
		users.Use(authMiddlewareHandler.AuthRequired())
		{
			users.GET("/", usersHandler.GetUserByUuid)
			users.PATCH("/", usersHandler.UpdateUser)
			users.DELETE("/", usersHandler.DeleteUser)
		}

		authItems := v1.Group("/items")
		authItems.Use(authMiddlewareHandler.AuthRequired())
		{
			authItems.POST("/", itemsHandler.AddItem)
			authItems.PATCH("/:uuid", itemsHandler.UpdateItem)
		}

		//authLists := v1.Group("/lists")
		//authLists.Use(authMiddlewareHandler.AuthRequired())
		//{
		//	authLists.POST("/", listHandler.AddList)
		//}

		//authListItems := v1.Group("/lists/:uuid")
		//authListItems.Use(authMiddlewareHandler.AuthRequired())
		//{
		//	authListItems.POST("/", listItemsHandler.AddItemToList)
		//}

		statuses := v1.Group("/statuses")
		statuses.Use(authMiddlewareHandler.AuthRequired())
		{
			statuses.POST("/", statusesHandler.AddStatus)
			statuses.GET("/", statusesHandler.GetStatusesForUser)
		}

		// Admin routes
		admin := v1.Group("/admin")
		admin.Use(authMiddlewareHandler.AuthRequired())
		admin.Use(superUserMiddlewareHandler.SuperUserStatusRequired())
		// TODO add middleware to check superuser status
		{
			admin.GET("/users", usersHandler.GetAllUsers)
			admin.DELETE("/users/:uuid", usersHandler.DeleteUser)
		}
	}
}
