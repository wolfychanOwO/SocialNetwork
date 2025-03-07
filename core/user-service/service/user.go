package service

import (
    "SocialNetwork/core/user-service/models"
    "SocialNetwork/core/user-service/repository"
)

type UserService struct {
    UserRepo *repository.UserRepository
}

func (s *UserService) CreateUser(user *models.User) error {
    return s.UserRepo.CreateUser(user)
}

func (s *UserService) AuthenticateUser(login, password string) (*models.User, error) {
    user, err := s.UserRepo.GetUserByLogin(login)
    if err != nil {
        return nil, err
    }

    // Здесь должна быть проверка пароля (например, с использованием bcrypt)
    if user.PasswordHash != password {
        return nil, err
    }

    return user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
    return s.UserRepo.UpdateUser(user)
}

func (s *UserService) GetUserByLogin(login string) (*models.User, error) {
    return s.UserRepo.GetUserByLogin(login)
}

