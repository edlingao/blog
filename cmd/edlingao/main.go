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
	configurator.CLIService()

	if len(arguments) == 0 || arguments[0] == "start" {
		log.Println("Starting server...")
		return
	}

	if arguments[0] == "save" {
		configurator.CLISave(arguments[1])
		return
	}
}
