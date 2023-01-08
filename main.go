package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vijeyash1/server/handlers"
	"github.com/vijeyash1/server/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}
	val := os.Getenv("AUTH0_CLIENT_ID")
	log.Println(val)
	auth, err := handlers.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	routes.InitRoutes(auth)

}
