package site

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/dbut2/home/pkg/log"
	"github.com/dbut2/home/pkg/pages"
	"github.com/go-chi/chi"
)

type Server struct{}

func (s *Server) Run(port string) {
	r := chi.NewRouter()

	r.Get("/404", func(w http.ResponseWriter, r *http.Request) {
		errorPage(w, r, 404, "Page not found 🔎")
	})

	r.Get("/500", func(w http.ResponseWriter, r *http.Request) {
		errorPage(w, r, 500, "Oops, server issue 🤪")
	})

	r.Handle("/static/*", http.FileServer(http.Dir("./web")))

	r.Get("/{page}", func(w http.ResponseWriter, r *http.Request) {
		page := chi.URLParam(r, "page")
		showPage(w, r, page)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		showPage(w, r, "home")
	})

	r.Get("/load", func(w http.ResponseWriter, r *http.Request) {
		client, err := datastore.NewClient(context.Background(), "dylanbutler")
		if err != nil {
			respondError(w, r, err)
			return
		}

		p := pages.Page{
			Title: "Home Page 🏡",
			Blocks: []pages.Block{
				{
					Type:    pages.Image,
					Index:   1,
					Content: "static/s1.jpg",
				},
				{
					Type:    pages.Text,
					Index:   2,
					Content: "This is some text",
				},
			},
		}

		err = p.Store(client, "home")
		if err != nil {
			respondError(w, r, err)
		}
	})

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		panic(err)
	}
}

func respondError(w http.ResponseWriter, r *http.Request, err error) {
	debug := true

	if debug {
		respondJSON(w, r, http.StatusInternalServerError, err.Error())
	} else {
		http.Redirect(w, r, "500", http.StatusTemporaryRedirect)
	}
}

func respondJSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	j, err := json.Marshal(v)
	if err != nil {
		respondError(w, r, err)
		return
	}
	w.WriteHeader(status)
	_, err = w.Write(j)
	if err != nil {
		respondError(w, r, err)
	}
}

func showPage(w http.ResponseWriter, r *http.Request, page string) {
	client, err := datastore.NewClient(context.Background(), "dylanbutler")
	if err != nil {
		respondError(w, r, err)
		return
	}
	p, err := pages.PageFromName(client, page)
	if err != nil {
		if err == pages.ErrNotFound {
			http.Redirect(w, r, "404", http.StatusTemporaryRedirect)
			return
		}
		respondError(w, r, err)
		return
	}
	err = p.Display(w)
	if err != nil {
		respondError(w, r, err)
		return
	}
}

func errorPage(w http.ResponseWriter, r *http.Request, status int, msg string) {
	p := pages.Page{
		Title: fmt.Sprintf("%d", status),
		Blocks: []pages.Block{
			{
				Type:    pages.Text,
				Index:   0,
				Content: fmt.Sprintf("Error %d: %s", status, msg),
			},
		},
	}

	err := p.Display(w)

	if err != nil {
		log.Error(err)
		respondError(w, r, err)
		return
	}
}
