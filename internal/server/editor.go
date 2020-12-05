package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dbut2/home/pkg/log"
	"github.com/dbut2/home/pkg/pages"
	"github.com/go-chi/chi"
)

func (s *Server) editor() chi.Router {
	editor := chi.NewRouter()

	editor.Get("/", func(w http.ResponseWriter, r *http.Request) {
		s.listPages(w, r)
	})

	editor.Get("/new", func(w http.ResponseWriter, r *http.Request) {
		s.newPage(w, r)
	})

	editor.Post("/new", func(w http.ResponseWriter, r *http.Request) {
		s.savePage(w, r)
	})

	editor.Get("/{page}", func(w http.ResponseWriter, r *http.Request) {
		page := chi.URLParam(r, "page")
		s.editPage(w, r, page)
	})

	editor.Post("/{page}", func(w http.ResponseWriter, r *http.Request) {
		page := chi.URLParam(r, "page")
		s.updatePage(w, r, page)
	})

	return editor
}

func (s *Server) newPage(w http.ResponseWriter, r *http.Request) {
	p := &pages.Page{}
	s.displayEdit(w, r, p)
}

func (s *Server) savePage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	name := r.FormValue("name")
	visible := r.FormValue("visible") == "on"
	title := r.FormValue("title")
	content := r.FormValue("content")
	p := pages.PageFromKey(pages.PageKey(name))
	p.Visible = visible
	p.Title = title
	p.Content = content
	err = p.Post(s.client)
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	http.Redirect(w, r, name, http.StatusSeeOther)
}

func (s *Server) editPage(w http.ResponseWriter, r *http.Request, page string) {
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
	s.displayEdit(w, r, p)
}

func (s *Server) updatePage(w http.ResponseWriter, r *http.Request, page string) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	name := r.FormValue("name")
	visible := r.FormValue("visible") == "on"
	title := r.FormValue("title")
	content := strings.ReplaceAll(r.FormValue("content"), "\r", "")
	pKey := pages.PageKey(page)
	p, err := pages.GetPageFromKey(s.client, pKey)
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	err = p.UpdateKey(s.client, pages.PageKey(name))
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	p.Visible = visible
	p.Title = title
	p.Content = content
	err = p.Put(s.client)
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	http.Redirect(w, r, name, http.StatusSeeOther)
}

func (s *Server) displayEdit(w http.ResponseWriter, r *http.Request, p *pages.Page) {
	err := PageDisplayEdit(p, w)
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
}

func (s *Server) listPages(w http.ResponseWriter, r *http.Request) {
	ps, err := pages.AllPages(s.client)
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
	p := &pages.Page{}
	shown := []*pages.Page{}
	hidden := []*pages.Page{}
	for _, pn := range ps {
		if pn.Visible {
			shown = append(shown, pn)
		} else {
			hidden = append(hidden, pn)
		}
	}

	content := "# Pages:\n"
	for _, pn := range shown {
		content += fmt.Sprintf("- [%s](%s) (%s) ([Edit](/edit/%s))\n", pn.Title, pn.Name(), pn.Name(), pn.Name())
	}
	content += "- [New Page](/edit/new)\n"
	content += "\n"
	content += "### Hidden:\n"
	for _, pn := range hidden {
		content += fmt.Sprintf("- **%s** (%s) ([Edit](/edit/%s))\n", pn.Title, pn.Name(), pn.Name())
	}
	p.Content = content
	err = PageDisplay(p, w)
	if err != nil {
		log.Error(err)
		s.respondError(w, r, err)
		return
	}
}
