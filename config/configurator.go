package config

import (
	"log"
	"os"

	"github.com/edlingao/internal/blog/adapters"
	"github.com/edlingao/internal/blog/ports"

	userAdapters "github.com/edlingao/internal/auth/adapters"
	userPorts "github.com/edlingao/internal/auth/ports"
	"github.com/edlingao/internal/pkg/database"
	"github.com/edlingao/web"
	"github.com/edlingao/web/template/layout"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Configurator struct {
	echo  *echo.Echo
	root  *echo.Group
	apiv1 *echo.Group
	db    *sqlx.DB
	CLI   ports.BlogCLISerivce
	Blog  ports.BlogService
	Users userPorts.UserService
}

func NewConfigurator() *Configurator {
	isDev := os.Getenv("ENV") != "prod"

	if !isDev {
		if err := layout.LoadManifest(); err != nil {
			log.Printf("Warning: Failed to load manifest: %v", err)
		}
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.StaticFS("/static", web.Static)

	root := e.Group("")
	apiv1 := root.Group("/api/v1")
	db := database.New()

	return &Configurator{
		echo:  e,
		root:  root,
		apiv1: apiv1,
		db:    db,
	}
}

func (configurator *Configurator) CLIService() *Configurator {
	blogRepo := adapters.NewBlogRepo(configurator.db)
	cliService := adapters.NewCLIService(
		blogRepo,
		configurator.Users,
	)

	configurator.CLI = cliService

	return configurator
}

func (configurator *Configurator) AddUserService() *Configurator {
	userRepo := userAdapters.NewUserRepo(configurator.db)
	userService := userAdapters.NewUserService(
		userRepo,
		configurator.root,
	)

	configurator.Users = userService

	return configurator
}

func (configurator *Configurator) AddBlogService() *Configurator {
	blogRepo := adapters.NewBlogRepo(configurator.db)
	blogService := adapters.NewBlogService(
		configurator.root,
		blogRepo,
	)

	configurator.Blog = blogService

	return configurator
}

func (configurator *Configurator) CLISave(title string) error {
	return configurator.CLI.SaveEntry(title)
}

func (configurator *Configurator) CLIAddUser(username, password, role string) error {
	return configurator.Users.Register(username, password, role)
}

func (configurator *Configurator) ServerStart() error {
	return configurator.echo.Start(":3011")
}
