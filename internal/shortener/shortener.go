package shortener

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/dbut2/home/pkg/url"
	"math/rand"
)

type Shortener struct {
	client *datastore.Client
}

func NewShortener() (*Shortener, error) {
	client, err := datastore.NewClient(context.Background(), "dylanbutler")
	if err != nil {
		return nil, err
	}

	return &Shortener{
		client: client,
	}, nil
}

func (s *Shortener) Shorten(url url.URL) (string, error) {
	code := fmt.Sprintf("%d", rand.Int())
	key := datastore.NameKey("urlcode", code, nil)
	_, err := s.client.Put(context.Background(), key, &url)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *Shortener) Lengthen(code string) (url.URL, error) {
	u := url.URL{}
	key := datastore.NameKey("urlcode", code, nil)
	err := s.client.Get(context.Background(), key, &u)
	if err != nil {
		return url.URL{}, err
	}
	return u, nil
}
