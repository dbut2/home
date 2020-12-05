package server

import (
	"context"

	"cloud.google.com/go/datastore"
)

type Config struct {
	Project string
	Admin   bool
}

func FromConfig(c *Config) (Server, error) {
	s := Server{
		Config: c,
	}
	client, err := datastore.NewClient(context.Background(), c.Project)
	if err != nil {
		return Server{}, err
	}
	s.client = client
	return s, nil
}
