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

	s := site.Server{}
	s.Run(port)
}
