package main

import (
	"log"
	"techzone/internal/notification/app"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(".env not found")
	}

	deps, err := app.BuildDependencies()
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.New(deps)
	if err != nil {
		log.Fatal(err)
	}

	defer application.Close()

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
