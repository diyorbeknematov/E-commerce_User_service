package postgres

import (
	"auth-service/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthenticatonRepository interface {
	Register(user models.Register) (*models.Response, error)
	Login(login models.LoginRequest) (*models.User, error)
	SaveRefreshToken(req models.TokenRequest) error
	InvalidateRefreshToken(username string) error
	IsRefreshTokenValid(username string) (bool, error)
}

type authenticateRepositoryImpl struct {
	db *sqlx.DB
}

func NewAuthenticatonRepository(db *sqlx.DB) AuthenticatonRepository {
	return &authenticateRepositoryImpl{db: db}
}

func (r *authenticateRepositoryImpl) Register(user models.Register) (*models.Response, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return nil, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	var id string
	err = tx.QueryRow(`
        INSERT INTO users (
            username,
            fullname,
            email,
            phone,
			image,
			role,
			password
        ) VALUES (
		 	$1, 
			$2, 
			$3, 
			'Unknown Phone', 
			'Unknown Image',
			'user', 
			$4) RETURNING id`,
		user.Username,
		user.FullName,
		user.Email,
		user.Password,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO user_locations (
			user_id,
			address,
			city,
			state,
			country,
			postal_code
		)
		VALUES (
			$1,
			'Unknown address',
			'Unknown city',
			'Unknown state',
            'Unknown country',
            'Unknown postal_code'
		)
	`, id)

	if err != nil {
		return nil, err
	}

	return &models.Response{
		Success: true,
		Message: "User registered successfully",
	}, nil
}

func (r *authenticateRepositoryImpl) Login(login models.LoginRequest) (*models.User, error) {
	var user models.User
	err := r.db.Get(&user, `
	SELECT 
		id, 
		username, 
		password,
		role
	FROM 
		users WHERE username=$1`,
		login.Username)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authenticateRepositoryImpl) SaveRefreshToken(req models.TokenRequest) error {
	_, err := r.db.Exec(`
		DELETE FROM 
			refresh_tokens 
		WHERE
			username = $1
	`, req.Username)

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
		fmt.Println(err)
		return err
	}

	_, err = r.db.Exec(`
        INSERT INTO refresh_tokens (
			username, 
			token,
			expires_at
		)
        VALUES ($1, 
				$2, 
				$3)
    `, req.Username, req.RefreshToken, req.ExpiresAt)

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
		fmt.Println(err)
	}

	return err
}

func (r *authenticateRepositoryImpl) InvalidateRefreshToken(username string) error {
	_, err := r.db.Exec(`
        DELETE FROM 
			refresh_tokens 
        WHERE 
            username=$1
    `, username)

	return err
}

func (r *authenticateRepositoryImpl) IsRefreshTokenValid(username string) (bool, error) {
	var count int
	err := r.db.Get(&count, `
        SELECT 
            count(*)
        FROM 
            refresh_tokens 
        WHERE 
            username=$1 AND 
            expires_at > CURRENT_TIMESTAMP
    `, username)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
