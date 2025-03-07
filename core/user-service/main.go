package main

import (
    "database/sql"
    "log"
    "net/http"
    "SocialNetwork/core/user-service/handlers"
    "SocialNetwork/core/user-service/repository"
    "SocialNetwork/core/user-service/service"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:userpass@172.20.0.2:5432/users?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    userRepo := &repository.UserRepository{DB: db}
    userService := &service.UserService{UserRepo: userRepo}
    userHandler := &handlers.UserHandler{UserService: userService}

    r := gin.Default()

    r.POST("/register", userHandler.RegisterUser)
    r.POST("/authenticate", userHandler.AuthenticateUser)
    r.PUT("/user", userHandler.UpdateUser)
    r.GET("/user/:login", userHandler.GetUserProfile)

    if err := http.ListenAndServe(":8081", r); err != nil {
        log.Fatal(err)
    }
}

