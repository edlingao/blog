package ports

import "github.com/edlingao/internal/auth/core"

type UserRepository interface {
	AddUser(user *core.User) (*core.User, error)
	GetUserByUsername(username string) (*core.User, error)
	UpdateUser(user *core.User) (*core.User, error)
}
