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
		rp.Select("https://weight-ehyj2hks2q-km.a.run.app", rp.PathIsAt("/weight"), rp.WithOIDC()),
		rp.Select("https://auth-ehyj2hks2q-km.a.run.app", rp.HostMatches("auth.dylanbutler.dev"), rp.WithOIDC()),
		rp.Select("https://dylanbutler-dev.web.app", rp.Always()),
	)

	err := http.ListenAndServe(":"+port, proxy)
	if err != nil {
		panic(err)
	}
}
