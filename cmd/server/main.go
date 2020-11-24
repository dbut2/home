package main

import (
	"os"

	"github.com/dbut2/home/internal/site"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	server := site.Server{}
	server.Run(port)
}
