package main

import (
	"net/http"
	"os"

	"github.com/dbut2/reverse-proxy"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	proxy := rp.New(
		rp.Select("https://shortener-prod-hqkniphctq-km.a.run.app/shorten", rp.PathIsAt("/shorten")),
		rp.Select("https://dylanbutler-dev.web.app", rp.Always()),
	)

	err := http.ListenAndServe(":"+port, proxy)
	if err != nil {
		panic(err)
	}
}
