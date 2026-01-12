package main

import (
	"log"
	"os"

	"github.com/edlingao/config"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	arguments := os.Args[1:]
	configurator := config.NewConfigurator()
	configurator.
		AddUserService().
		CLIService().
		AddBlogService()

	if len(arguments) == 0 || arguments[0] == "start" {
		log.Fatal(configurator.ServerStart())
		return
	}

	if arguments[0] == "save" {
		configurator.CLISave(arguments[1])
		return
	}

	if arguments[0] == "adduser" {
		configurator.CLIAddUser(arguments[1], arguments[2], arguments[3])
	}

}
