package service

import (
    "testing"
    "time"
    "SocialNetwork/core/user-service/models"
    "SocialNetwork/core/user-service/repository"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockUserRepository) GetUserByLogin(login string) (*models.User, error) {
    args := m.Called(login)
    return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user *models.User) error {
    args := m.Called(user)
    return args.Error(0)
}

func TestCreateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := &UserService{UserRepo: mockRepo}

    user := &models.User{
        Login:        "user1",
        PasswordHash: "hash",
        Email:        "user1@example.com",
        FirstName:    "John",
        LastName:     "Doe",
        BirthDate:    "1990-01-01",
        PhoneNumber:  "1234567890",
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    mockRepo.On("CreateUser", user).Return(nil)

    err := service.CreateUser(user)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

func TestAuthenticateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := &UserService{UserRepo: mockRepo}

    user := &models.User{
        Login:        "user1",
        PasswordHash: "hash",
    }

    mockRepo.On("GetUserByLogin", "user1").Return(user, nil)

    authenticatedUser, err := service.AuthenticateUser("user1", "hash")
    assert.NoError(t, err)
    assert.Equal(t, "user1", authenticatedUser.Login)
    mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
    mockRepo := new(MockUserRepository)
    service := &UserService{UserRepo: mockRepo}

    user := &models.User{
        ID:          1,
        Email:       "newemail@example.com",
        FirstName:   "John",
        LastName:    "Smith",
        BirthDate:   "1990-01-01",
        PhoneNumber: "0987654321",
        UpdatedAt:   time.Now(),
    }

    mockRepo.On("UpdateUser", user).Return(nil)

    err := service.UpdateUser(user)
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}

