package config

import (
	"github.com/edlingao/internal/blog/adapters"
	"github.com/edlingao/internal/blog/ports"
	"github.com/edlingao/internal/pkg/database"
	"github.com/jmoiron/sqlx"
)

type Configurator struct {
	db  *sqlx.DB
	CLI ports.BlogCLISerivce
}

func NewConfigurator() *Configurator {
	db := database.New()
	return &Configurator{
		db: db,
	}
}

func (configurator *Configurator) CLIService() *Configurator {
	blogRepo := adapters.NewBlogRepo(configurator.db)
	cliService := adapters.NewCLIService(
		blogRepo,
	)

	configurator.CLI = cliService

	return configurator
}


func (configurator *Configurator) CLISave(title string) error {
	return configurator.CLI.SaveEntry(title)
}
