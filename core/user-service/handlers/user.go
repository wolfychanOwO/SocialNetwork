package handlers

import (
    "net/http"
    "SocialNetwork/core/user-service/models"
    "SocialNetwork/core/user-service/service"
    "github.com/gin-gonic/gin"
)

type UserHandler struct {
    UserService *service.UserService
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.UserService.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) AuthenticateUser(c *gin.Context) {
    var auth struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }
    if err := c.ShouldBindJSON(&auth); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.UserService.AuthenticateUser(auth.Login, auth.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.UserService.UpdateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserProfile(c *gin.Context) {
    login := c.Param("login")
    user, err := h.UserService.GetUserByLogin(login)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

