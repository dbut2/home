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

	selectors := []rp.Selector{
		rp.Select("https://shortener-prod-hqkniphctq-km.a.run.app/shorten", rp.PathRule("/shorten")),
		rp.Select("https://dylanbutler-dev.web.app", rp.BaseRule()),
	}
	proxy := rp.New(selectors...)

	err := http.ListenAndServe(":"+port, proxy)
	if err != nil {
		panic(err)
	}
}
