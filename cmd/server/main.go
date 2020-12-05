package main

import (
	"os"

	"github.com/dbut2/home/internal/server"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	s, _ := server.FromConfig(&server.Config{
		Project: "dylanbutler",
		Admin:   true,
	})
	s.Run(port)
}
