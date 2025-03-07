package handlers

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"SocialNetwork/core/user-service/handlers"
    "SocialNetwork/core/user-service/models"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) CreateUser(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserService) AuthenticateUser(login, password string) (*models.User, error) {
    args := m.Called(login, password)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserService) GetUserByLogin(login string) (*models.User, error) {
    args := m.Called(login)
    return args.Get(0).(*models.User), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
    mockService := new(MockUserService)
    handler := &handlers.UserHandler{UserService: mockService}

    user := &models.User{
        Login:        "user1",
        PasswordHash: "hash",
        Email:        "user1@example.com",
        FirstName:    "John",
        LastName:     "Doe",
        BirthDate:    "1990-01-01",
        PhoneNumber:  "1234567890",
    }

    mockService.On("CreateUser", user).Return(nil)

    r := gin.Default()
    r.POST("/register", handler.RegisterUser)

    body, _ := json.Marshal(user)
    req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}

