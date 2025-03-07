package repository

import (
    "testing"
    "time"
	"SocialNetwork/core/user-service/models"
	"SocialNetwork/core/user-service/repository"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := &repository.UserRepository{DB: db}

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

    mock.ExpectQuery("INSERT INTO users").
        WithArgs(user.Login, user.PasswordHash, user.Email, user.FirstName, user.LastName, user.BirthDate, user.PhoneNumber, user.CreatedAt, user.UpdatedAt).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    err = repo.CreateUser(user)
    assert.NoError(t, err)
    assert.Equal(t, 1, user.ID)
}

func TestGetUserByLogin(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := &repository.UserRepository{DB: db}

    rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "email", "first_name", "last_name", "birth_date", "phone_number", "created_at", "updated_at"}).
        AddRow(1, "user1", "hash", "user1@example.com", "John", "Doe", "1990-01-01", "1234567890", time.Now(), time.Now())

    mock.ExpectQuery("SELECT id, login, password_hash, email, first_name, last_name, birth_date, phone_number, created_at, updated_at FROM users WHERE login = ?").
        WithArgs("user1").
        WillReturnRows(rows)

    user, err := repo.GetUserByLogin("user1")
    assert.NoError(t, err)
    assert.Equal(t, "user1", user.Login)
}

func TestUpdateUser(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := &repository.UserRepository{DB: db}

    user := &models.User{
        ID:          1,
        Email:       "newemail@example.com",
        FirstName:   "John",
        LastName:    "Smith",
        BirthDate:   "1990-01-01",
        PhoneNumber: "0987654321",
        UpdatedAt:   time.Now(),
    }

    mock.ExpectExec("UPDATE users SET email = ?, first_name = ?, last_name = ?, birth_date = ?, phone_number = ?, updated_at = ? WHERE id = ?").
        WithArgs(user.Email, user.FirstName, user.LastName, user.BirthDate, user.PhoneNumber, user.UpdatedAt, user.ID).
        WillReturnResult(sqlmock.NewResult(1, 1))

    err = repo.UpdateUser(user)
    assert.NoError(t, err)
}

