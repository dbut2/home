package pages

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"cloud.google.com/go/datastore"
	"github.com/dbut2/home/pkg/log"
)

var (
	ErrNotFound = errors.New("page not found")
)

type Page struct {
	Title  string
	Blocks []Block `datastore:"-"`
}

func PageFromName(client *datastore.Client, name string) (Page, error) {
	p := Page{}
	key := datastore.NameKey("Page", name, nil)
	err := client.Get(context.Background(), key, &p)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return Page{}, ErrNotFound
		}
		log.Error(err)
		return Page{}, err
	}
	_, err = client.GetAll(context.Background(), datastore.NewQuery("Block").Ancestor(key).Order("Index"), &p.Blocks)
	if err != nil {
		log.Error(err)
		return Page{}, err
	}
	return p, nil
}

func (p Page) Display(w http.ResponseWriter) error {
	t, err := template.New("Page").Parse(page)
	if err != nil {
		return err
	}
	err = t.Execute(w, p)
	if err != nil {
		return err
	}
	return nil
}

func (p Page) Store(client *datastore.Client, name string) error {
	key := datastore.NameKey("Page", name, nil)
	_, err := client.Put(context.Background(), key, &p)
	if err != nil {
		return err
	}
	for _, b := range p.Blocks {
		err = b.Store(client, key)
		if err != nil {
			return err
		}
	}
	return nil
}

type Block struct {
	Type    BlockType
	Index   int
	Content string
}

func (b Block) Display() template.HTML {
	switch b.Type {
	case Image:
		return template.HTML(fmt.Sprintf("<img src=\"static/%s\" />", b.Content))
	case Text:
		return template.HTML(fmt.Sprintf("<p>%s</p>", b.Content))
	default:
		return ""
	}
}

func (b Block) Store(client *datastore.Client, parent *datastore.Key) error {
	fmt.Println(int64(b.Index))
	key := datastore.IDKey("Block", int64(b.Index), parent)
	_, err := client.Put(context.Background(), key, &b)
	if err != nil {
		return err
	}
	return nil
}

type BlockType int

const (
	Undefined BlockType = iota

	HTML
	Text
	Image
)

const (
	page = `<!DOCTYPE html>
	<html lang="en">
	<head>
	<meta charset="utf-8" />
	<title>{{ .Title }}</title>
	<link rel="stylesheet" href="static/style.css" />
	</head>
	<body>
	<div id="container">
	{{ range .Blocks }}
	<div>{{ .Display }}</div>
	{{ end }}
	</div>
	</body>
	</html>`
)
