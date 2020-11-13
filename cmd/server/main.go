package main

import (
	"github.com/dbut2/home/internal/site"
	"os"
)

func main() {
	port := os.Getenv("port")

	if port == "" {
		port = "8080"
	}
	server := site.Server{}
	server.Run(port)
}
