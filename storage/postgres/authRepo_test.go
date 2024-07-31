package postgres

import (
	"auth-service/models"
	"context"
	"log"
	"testing"
)

func TestConnDB(t *testing.T) {
	db, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err != nil {
		t.Fatalf("Error in connect database: %v", err)
	}
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}
}

func TestRegister(t *testing.T) {
	conn, err := Connection()
	if err != nil {
		log.Fatal(err)
	}

	tr, err := conn.BeginTx(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tr.Rollback()

	auth := authenticateRepositoryImpl{conn}
	user := models.Register{
		Username: "test",
		Password: "test",
		Email:    "test@test.com",
		FullName: "Kimdur Kimdurov",
	}
	result, err := auth.Register(user)
	if err != nil {
		log.Fatal(err)
	}
	if !result.Success {
		log.Fatal(result.Message)
	}
	_, err = conn.Exec("delete from users where username = $1", user.Username)
	if err != nil {
		log.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	conn, err := Connection()
	if err != nil {
		log.Fatal(err)
	}

	auth := authenticateRepositoryImpl{conn}

	_, err = conn.Exec("insert into users(username, password, email) values($1, $2, $3)",
		"test", "test", "test@test.com")
	if err != nil {
		log.Fatal(err)
	}

	login := models.LoginRequest{
		Username: "test",
	}
	result, err := auth.Login(login)
	if err != nil {
		log.Fatal(err)
	}
	
	if result.Username!= "test" {
        t.Error("Test failed, expected 'test' but got", result.Username)
    }

	_, err = conn.Exec("DELETE FROM users WHERE username = $1", "test")
	if err != nil {
		log.Fatal(err)
	}
}

func TestSaveRefreshToken(t *testing.T) {
	conn, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	auth := authenticateRepositoryImpl{conn}

	token := models.TokenRequest{
		Username:     "test",
		RefreshToken: "test",
		ExpiresAt:    "2024-12-05",
	}
	err = auth.SaveRefreshToken(token)
	if err != nil {
		log.Fatal(err)
	}

	conn.Exec("delete from users where username = $1", "test")
}

func TestInvalidateRefreshToken(t *testing.T) {
	conn, err := Connection()
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Exec("insert into refresh_tokens(username, token, expires_at) values($1, $2, $3)",
		"test", "test", "2024-12-05")
	if err != nil {
		t.Fatalf("Error inserting refresh token: %v", err)
	}

	auth := authenticateRepositoryImpl{conn}

	err = auth.InvalidateRefreshToken("est")
	if err != nil {
		log.Fatal(err)
	}
}

func TestIsValidRefreshToken(t *testing.T) {
	conn, err := Connection()
	if err != nil {
		log.Fatal(err)
	}
	auth := authenticateRepositoryImpl{conn}

	token := models.TokenRequest{
		Username:     "test",
		RefreshToken: "test",
		ExpiresAt:    "2024-12-05",
	}
	err = auth.SaveRefreshToken(token)
	if err != nil {
		log.Fatal(err)
	}

	check, err := auth.IsRefreshTokenValid("test")
	if err != nil {
		log.Fatal(err)
	}
	if !check {
		log.Fatal("refresh token is not valid")
	}

	_, err = conn.Exec("DELETE FROM refresh_tokens WHERE token = $1", token.RefreshToken)
	if err != nil {
		log.Fatal(err)
	}
}
