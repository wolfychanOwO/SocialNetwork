package models

import "time"

type User struct {
    ID           int       `json:"id"`
    Login        string    `json:"login"`
    PasswordHash string    `json:"-"`
    Email        string    `json:"email"`
    FirstName    string    `json:"first_name"`
    LastName     string    `json:"last_name"`
    BirthDate    string    `json:"birth_date"`
    PhoneNumber  string    `json:"phone_number"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

