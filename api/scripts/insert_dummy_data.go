package main

import (
	"codeberg.org/sporiff/eigakanban/config"
	queries "codeberg.org/sporiff/eigakanban/db/sqlc"
	"codeberg.org/sporiff/eigakanban/services"
	"codeberg.org/sporiff/eigakanban/types"
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"log"
)

func main() {
	dbConfig := config.LoadDBConfig()

	db, err := config.ConnectDB(dbConfig)
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %v", err)
	}
	defer db.Close()

	q := queries.New(db)

	authService := services.NewAuthService(q)
	usersService := services.NewUsersService(q)
	itemService := services.NewItemsService(q)

	createDummyUsers(context.Background(), authService, usersService)
	createDummyItems(context.Background(), itemService)

	log.Println("Dummy data generation complete!")
}

func createDummyUsers(ctx context.Context, authService *services.AuthService, usersService *services.UsersService) {
	total, err := usersService.GetUserCount(ctx)
	if err != nil {
		log.Fatalf("Couldn't get all users: %v", err)
	}

	if total > 0 {
		log.Printf("%d users already present", total)
		return
	}

	for i := 0; i < 10; i++ {
		user := types.RegisterUserRequest{
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: gofakeit.Password(true, true, true, true, false, 10),
		}
		_, err := authService.RegisterUser(ctx, user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", user.Username, err)
		} else {
			log.Printf("Created user: %s", user.Username)
		}
	}
}

func createDummyItems(ctx context.Context, itemService *services.ItemsService) {
	total, err := itemService.GetItemsCount(ctx)
	if err != nil {
		log.Fatalf("Couldn't get all items: %v", err)
	}

	if total > 0 {
		log.Printf("%d items already present", total)
		return
	}

	for i := 0; i < 20; i++ {
		item := types.AddItemRequest{
			ItemTitle: gofakeit.MovieName(),
		}
		_, err := itemService.AddItem(ctx, item)
		if err != nil {
			log.Printf("Failed to create item %s: %v", item.ItemTitle, err)
		} else {
			log.Printf("Created item: %s", item.ItemTitle)
		}
	}
}
