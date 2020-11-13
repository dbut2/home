package main

import (
	"fmt"
	"github.com/dbut2/home/internal/shortener"
	"github.com/dbut2/home/internal/site"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	s, err := shortener.NewShortener()
	if err != nil {
		fmt.Println(err)
	}
	server := site.Shortener{Shortener: s}
	server.Run(port)
}
