package db

import (
	"bpzh-api/internal/model"
	"cloud.google.com/go/firestore"
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

func (r *Repo) GetUserByVkId(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	iter := r.FS.Collection("users").Where("vkid", "==", id).Documents(ctx)
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

func (r *Repo) UpdateUserOnVk(ctx context.Context, docId string, u model.User) error {
	_, err := r.FS.Collection("users").Doc(docId).Update(ctx, []firestore.Update{
		{
			Path:  "firstname",
			Value: u.FirstName,
		},
		{
			Path:  "lastname",
			Value: u.LastName,
		},
		{
			Path:  "vkdomain",
			Value: u.VkDomain,
		},
		{
			Path:  "vkid",
			Value: u.VkId,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) CreateUserOnVk(ctx context.Context, u model.User) error {
	_, _, err := r.FS.Collection("users").Add(ctx, map[string]interface{}{
		"firstname": u.FirstName,
		"lastname":  u.LastName,
		"vkdomain":  u.VkDomain,
		"vkid":      u.VkId,
	})
	if err != nil {
		return err
	}

	return nil
}
