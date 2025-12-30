package repositories

import (
	"github.com/edlingao/internal/auth/core"
	"github.com/edlingao/internal/auth/queries"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (userRepo *UserRepo) AddUser(user *core.User) (*core.User, error) {
	uuid := uuid.NewString()
	user.ID = uuid
	_, err := userRepo.db.NamedExec(queries.AddUserQuery, user)
	if err != nil {
		return &core.User{}, err
	}

	return userRepo.GetUserByUsername(user.Username)
}

func (userRepo *UserRepo) GetUserByUsername(username string) (*core.User, error) {
	user := core.NewUser(username, "")
	err := userRepo.db.Get(user, queries.GetUserByUsernameQuery, username)
	if err != nil {
		return &core.User{}, err
	}

	return user, nil
}

func (userRepo *UserRepo) UpdateUser(user *core.User) (*core.User, error) {
	_, err := userRepo.db.NamedExec(queries.UpdateUserQuery, user)
	if err != nil {
		return &core.User{}, err
	}

	return userRepo.GetUserByUsername(user.Username)
}
