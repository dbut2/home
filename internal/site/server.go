package site

import (
	"encoding/json"
	"fmt"
	"github.com/dbut2/home/internal/shortener"
	"github.com/dbut2/home/pkg/url"
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

type Shortener struct {
	*shortener.Shortener
}

func (s *Shortener) Run(port string) {
	r := chi.NewRouter()

	r.Get("/*", s.lengthen)
	r.Post("/shorten", s.shorten)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}

func (s *Shortener) lengthen(w http.ResponseWriter, r *http.Request) {

	code := strings.TrimPrefix(r.URL.Path, "/")

	link, err := s.Lengthen(code)
	if err != nil {
		http.Redirect(w, r, "404", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, link.String(), http.StatusMovedPermanently)

}

func (s *Shortener) shorten(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		respondError(w, err)
		return
	}
	link := r.FormValue("url")
	code := r.FormValue("code")
	if link == "" {
		respondError(w, fmt.Errorf("url empty"))
		return
	}
	u, err := url.New(r.FormValue("url"))
	if err != nil {
		respondError(w, err)
		return
	}
	code, err = s.Shorten(u, code)
	if err != nil {
		respondError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, struct {
		URL  string
		Code string
	}{
		URL:  link,
		Code: code,
	})
}

type Server struct{}

func (s *Server) Run(port string) {
	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}

func respondError(w http.ResponseWriter, err error) {
	respondJSON(w, http.StatusInternalServerError, err)
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(status)
	_, err = w.Write(j)
	if err != nil {
		respondError(w, err)
	}
}
