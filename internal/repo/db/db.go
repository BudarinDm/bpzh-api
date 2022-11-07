package db

import (
	"bpzh-api/internal/config"
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
)

type Repo struct {
	FS *firestore.Client
}

func NewRepo(FS *firestore.Client) *Repo {
	return &Repo{
		FS: FS,
	}
}

func CreateFSConnections(cfg *config.DBConfig) (*firestore.Client, error) {
	err := ioutil.WriteFile("fs.json", []byte(cfg.FSConf), 0777)
	if err != nil {
		fmt.Println(err)
	}

	options := option.WithCredentialsFile("./fs.json")
	client, err := firestore.NewClient(context.Background(), "bpzh-info", options)
	if err != nil {
		return nil, err
	}
	return client, nil
}
