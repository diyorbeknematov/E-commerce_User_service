package postgres

import "github.com/jmoiron/sqlx"

type MainRepository interface {
	Authentications() AuthenticatonRepository
	Users() UserRepository
}

type mainRepositoryImpl struct {
	Authetication AuthenticatonRepository
	User          UserRepository
}

func NewMainRepository(db *sqlx.DB) MainRepository {
	return &mainRepositoryImpl{
		Authetication: NewAuthenticatonRepository(db),
		User:          NewUserRepository(db),
	}
}

func (r *mainRepositoryImpl) Authentications() AuthenticatonRepository {
	return r.Authetication
}

func (r *mainRepositoryImpl) Users() UserRepository {
	return r.User
}
