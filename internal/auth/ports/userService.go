package ports

import (
	"github.com/edlingao/internal/auth/core"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	LoginHandler(echo.Context) error
	LoginPostHandler(echo.Context) error
	Register(username, password, role string) error
	Login(username, password string) ( *core.User, error )
}
