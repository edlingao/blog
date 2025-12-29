package adapters

import (
	"log"

	"github.com/edlingao/internal/auth/core"
	"github.com/edlingao/internal/auth/ports"
	authLibrary "github.com/edlingao/internal/pkg/auth"
	"github.com/edlingao/internal/pkg/web"
	"github.com/edlingao/web/template/auth"
	"github.com/labstack/echo/v4"
)

type UserService struct {
	userRepo ports.UserRepository
	root     *echo.Group
}

func NewUserService(
	userRepo ports.UserRepository,
	root *echo.Group,
) *UserService {
	userService := &UserService{
		userRepo: userRepo,
		root:     root,
	}

	userService.root.GET("/login", userService.LoginHandler)
	userService.root.POST("/login", userService.LoginPostHandler)
	userService.root.GET("/logout", userService.LogoutHandler)

	return userService
}

func (userService *UserService) LoginHandler(c echo.Context) error {
	return web.Render(c, auth.Login(auth.VMLogin{}), 200)
}

func (userService *UserService) LoginPostHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := userService.Login(username, password)
	if err != nil {
		return web.Render(c, auth.Login(auth.VMLogin{
			Error: "Invalid username or password",
		}), 200)
	}

	// Set user session or token here
	err = authLibrary.SetAuthCookie(user.GenerateToken(), c)
	if err != nil {
		log.Println("Error setting auth cookie: ", err)
		return web.Render(c, auth.Login(auth.VMLogin{
			Error: "Internal server error",
		}), 200)
	}

	// For simplicity, we just redirect to home
	return c.Redirect(302, "/")
}

func (userService *UserService) LogoutHandler(c echo.Context) error {
	err := authLibrary.LogoutCookie(c)
	if err != nil {
		log.Println("Error clearing auth cookie: ", err)
	}

	return c.Redirect(302, "/login")
}

func (userService *UserService) Register(username, password, role string) error {
	user := core.NewUser(username, role)
	user.NewPassword("", password)
	user, err := userService.userRepo.AddUser(user)
	if err != nil {
		log.Println("Error registering user: ", err)
		return err
	}

	return nil
}

func (userService *UserService) Login(username, password string) (*core.User, error) {
	user, err := userService.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Println("Error logging in user: ", err)
		return &core.User{}, err
	}

	if !user.ValidatePassword(password) {
		log.Println("Invalid password for user: ", username)
		return &core.User{}, echo.ErrUnauthorized
	}

	return user, nil
}
