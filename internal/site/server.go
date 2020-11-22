package site

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type Server struct{}

func (s *Server) Run(port string) {
	r := chi.NewRouter()

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Handle("/static/*", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}

func respondError(w http.ResponseWriter, err error) {
	respondJSON(w, http.StatusInternalServerError, err.Error())
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
