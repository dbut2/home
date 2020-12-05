package pages

import (
	"context"
	"errors"
	"html/template"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/dbut2/home/pkg/log"
	"github.com/gomarkdown/markdown"
)

var (
	ErrNotFound = errors.New("page not found")
)

type Page struct {
	key     *datastore.Key
	Visible bool
	Title   string
	Content string
}

func PageKey(name string) *datastore.Key {
	return datastore.NameKey("Page", name, nil)
}

func PageFromKey(key *datastore.Key) *Page {
	return &Page{
		key: key,
	}
}

func GetPageFromKey(client *datastore.Client, key *datastore.Key) (*Page, error) {
	p := PageFromKey(key)
	err := p.Get(client)
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, ErrNotFound
		}
		log.Error(err)
		return nil, err
	}
	return p, nil
}

func AllPages(client *datastore.Client) ([]*Page, error) {
	keys, err := client.GetAll(context.Background(), datastore.NewQuery("Page"), &[]*Page{})
	if err != nil {
		return nil, err
	}
	pages := []*Page{}
	for _, key := range keys {
		p, err := GetPageFromKey(client, key)
		if err != nil {
			return nil, err
		}
		pages = append(pages, p)
	}
	return pages, nil
}

func (p *Page) Name() string {
	if p.key == nil {
		return ""
	}
	return p.key.Name
}

func (p *Page) ParseContent() template.HTML {
	md := []byte(p.Content)
	parsed := markdown.ToHTML(md, nil, nil)
	return template.HTML(strings.ReplaceAll(string(parsed), "<hr>", "</div></div><div><div>"))
}

func (p *Page) Get(client *datastore.Client) error {
	err := client.Get(context.Background(), p.key, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) Post(client *datastore.Client) error {
	err := client.Get(context.Background(), p.key, &Page{})
	if err != datastore.ErrNoSuchEntity {
		if err != nil {
			log.Error(err)
			return err
		}
		return errors.New("page name already exist")
	}
	_, err = client.Put(context.Background(), p.key, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) Put(client *datastore.Client) error {
	err := client.Get(context.Background(), p.key, &Page{})
	if err != nil {
		if err == datastore.ErrNoSuchEntity {
			return errors.New("page name not exist")
		}
		log.Error(err)
		return err
	}
	_, err = client.Put(context.Background(), p.key, p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) Delete(client *datastore.Client) error {
	err := client.Delete(context.Background(), p.key)
	if err != nil {
		return err
	}
	return nil
}

func (p *Page) UpdateKey(client *datastore.Client, key *datastore.Key) error {
	if p.key.Name == key.Name {
		return nil
	}
	err := client.Get(context.Background(), key, &Page{})
	if err != datastore.ErrNoSuchEntity {
		if err != nil {
			log.Error(err)
			return err
		}
		return errors.New("page name already exist")
	}
	t, err := client.NewTransaction(context.Background())
	if err != nil {
		return err
	}
	err = t.Delete(p.key)
	if err != nil {
		err = t.Rollback()
		return err
	}
	_, err = t.Put(key, p)
	if err != nil {
		err = t.Rollback()
		return err
	}
	_, err = t.Commit()
	if err != nil {
		return err
	}
	p.key = key
	return nil
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
	return p, nil
}

func (p Page) Store(client *datastore.Client, name string) error {
	key := datastore.NameKey("Page", name, nil)
	_, err := client.Put(context.Background(), key, &p)
	if err != nil {
		return err
	}
	return nil
}
