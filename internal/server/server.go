package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/dbut2/home/pkg/log"
	"github.com/dbut2/home/pkg/pages"
	"github.com/go-chi/chi"
)

type Server struct {
	Config *Config
	client *datastore.Client
}

func (s *Server) Run(port string) {
	r := chi.NewRouter()

	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		s.respondJSON(w, r, 200, "favicon")
	})

	r.Get("/404", func(w http.ResponseWriter, r *http.Request) {
		s.errorPage(w, r, 404, "Page not found 🔎")
	})

	r.Get("/500", func(w http.ResponseWriter, r *http.Request) {
		s.errorPage(w, r, 500, "Oops, server issue 🤪")
	})

	r.Handle("/static/*", http.FileServer(http.Dir("./web")))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s.showPage(w, r, "home")
	})

	if s.Config.Admin {
		r.Mount("/edit", s.editor())
	}

	r.Get("/{page}", func(w http.ResponseWriter, r *http.Request) {
		page := chi.URLParam(r, "page")
		s.showPage(w, r, page)
	})

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}

func (s *Server) respondError(w http.ResponseWriter, r *http.Request, err error) {
	if s.Config.Admin {
		s.respondJSON(w, r, http.StatusInternalServerError, err.Error())
	} else {
		http.Redirect(w, r, "/500", http.StatusTemporaryRedirect)
	}
}

func (s *Server) respondJSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		s.respondError(w, r, err)
		return
	}
	w.WriteHeader(status)
	_, err = w.Write(j)
	if err != nil {
		s.respondError(w, r, err)
	}
}

func (s *Server) showPage(w http.ResponseWriter, r *http.Request, page string) {
	p, err := pages.GetPageFromKey(s.client, pages.PageKey(page))
	if err != nil {
		if err == pages.ErrNotFound {
			http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
			return
		}
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	if !p.Visible {
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
		return
	}
	err = PageDisplay(p, w)
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
}

func (s *Server) errorPage(w http.ResponseWriter, r *http.Request, status int, msg string) {
	p := &pages.Page{
		Title:   fmt.Sprintf("%d", status),
		Content: fmt.Sprintf("Error %d: %s", status, msg),
	}

	err := PageDisplay(p, w)

	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
}
