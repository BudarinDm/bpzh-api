package db

import (
	"bpzh-api/internal/model"
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"time"
)

func (r *Repo) GetCodeInfo(ctx context.Context, id string) (*model.CodeInfo, error) {
	var tInfo model.CodeInfo
	dsnap, err := r.FS.Collection("auth").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	err = dsnap.DataTo(&tInfo)
	if err != nil {
		return nil, err
	}
	return &tInfo, nil
}

func (r *Repo) UpdateCode(ctx context.Context, docId string, code int64) error {
	_, err := r.FS.Collection("auth").Doc(docId).Set(ctx, map[string]interface{}{
		"code":         code,
		"send_code_at": time.Now(),
		"try_count":    5,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) UpdateTryCountCode(ctx context.Context, docId string, tryCnt int64) error {
	_, err := r.FS.Collection("auth").Doc(docId).Update(ctx, []firestore.Update{
		{
			Path:  "try_count",
			Value: tryCnt - 1,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) CreateToken(ctx context.Context, docId string, token, userDocId string, vkId int64) error {
	_, err := r.FS.Collection("auth").Doc(docId).Set(ctx, map[string]interface{}{
		"token":        token,
		"create_token": time.Now(),
		"vkid":         vkId,
		"id":           userDocId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) GetSessionByToken(ctx context.Context, token string) (*model.TokenInfo, error) {
	var tInfo model.TokenInfo
	iter := r.FS.Collection("auth").Where("token", "==", token).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		err = doc.DataTo(&tInfo)
		if err != nil {
			return nil, err
		}
	}
	return &tInfo, nil
}
