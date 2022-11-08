package db

import (
	"bpzh-api/internal/model"
	"context"
	"google.golang.org/api/iterator"
)

func (r *Repo) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	iter := r.FS.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var user model.User
		err = doc.DataTo(&user)
		if err != nil {
			return nil, err
		}

		user.DocId = doc.Ref.ID
		users = append(users, user)
	}

	return users, nil
}

func (r *Repo) GetUserByVkDomain(ctx context.Context, domain string) (*model.User, error) {
	var user model.User
	iter := r.FS.Collection("users").Where("vkdomain", "==", domain).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		err = doc.DataTo(&user)
		if err != nil {
			return nil, err
		}
		user.DocId = doc.Ref.ID
	}

	return &user, nil
}
