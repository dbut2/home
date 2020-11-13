package site

import (
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
		fmt.Println(err)
		return
	}
	link, err := url.New(r.FormValue("url"))
	if err != nil {
		fmt.Println(err)
		return
	}
	code, err := s.Shorten(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintf(w, "Code: %s, Url: %s", code, link.String())
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Server struct {}

func (s *Server) Run(port string) {
	r := chi.NewRouter()
	r.Handle("/*", http.FileServer(http.Dir("./web")))
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
