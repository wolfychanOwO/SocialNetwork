package repository

import (
    "database/sql"
    "time"
    "SocialNetwork/core/user-service/models"
)

type UserRepository struct {
    DB *sql.DB
}

func (r *UserRepository) CreateUser(user *models.User) error {
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    query := `INSERT INTO users (login, password_hash, email, first_name, last_name, birth_date, phone_number, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
    return r.DB.QueryRow(query, user.Login, user.PasswordHash, user.Email, user.FirstName, user.LastName, user.BirthDate, user.PhoneNumber, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
}

func (r *UserRepository) GetUserByLogin(login string) (*models.User, error) {
    user := &models.User{}
    query := `SELECT id, login, password_hash, email, first_name, last_name, birth_date, phone_number, created_at, updated_at 
              FROM users WHERE login = $1`
    err := r.DB.QueryRow(query, login).Scan(&user.ID, &user.Login, &user.PasswordHash, &user.Email, &user.FirstName, &user.LastName, &user.BirthDate, &user.PhoneNumber, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (r *UserRepository) UpdateUser(user *models.User) error {
    user.UpdatedAt = time.Now()
    query := `UPDATE users SET email = $1, first_name = $2, last_name = $3, birth_date = $4, phone_number = $5, updated_at = $6 
              WHERE id = $7`
    _, err := r.DB.Exec(query, user.Email, user.FirstName, user.LastName, user.BirthDate, user.PhoneNumber, user.UpdatedAt, user.ID)
    return err
}

