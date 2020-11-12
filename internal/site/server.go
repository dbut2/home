package site

import (
	"github.com/go-chi/chi"
	"net/http"
)

type Server struct {
}

func (s *Server) Run(port string) {
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Handle("/static/*", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}
