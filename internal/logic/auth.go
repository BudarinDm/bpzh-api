package logic

import (
	"bpzh-api/internal/model"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

func (l *Logic) GetLastSend(ctx context.Context, request model.RequestCode) (int64, error) {
	info, err := l.repo.GetCodeInfo(ctx, fmt.Sprintf("%s.%s", request.Scope, request.Login))
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return 0, nil
		}
		return 0, err
	}

	return info.SendCodeAt.Unix(), nil
}

func (l *Logic) UpdateCode(ctx context.Context, docId string, code int64) error {
	return l.repo.UpdateCode(ctx, docId, code)
}

func (l *Logic) CreateToken(ctx context.Context, docId, token, userDocId string, vkId int64) error {
	return l.repo.CreateToken(ctx, docId, token, userDocId, vkId)
}

func (l *Logic) SendCode(code, userId int) error {
	err := l.vkApiClient.SendCode(code, userId)
	if err != nil {
		return err
	}
	return nil
}

func (l *Logic) GetCodeCheck(ctx context.Context, scope, login string) (int64, int64, error) {
	info, err := l.repo.GetCodeInfo(ctx, fmt.Sprintf("%s.%s", scope, login))
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	err = l.repo.UpdateTryCountCode(ctx, fmt.Sprintf("%s.%s", scope, login), info.TryCount)
	if err != nil {
		return 0, 0, err
	}

	return info.Code, info.TryCount, nil
}

func (l *Logic) GetSessionByToken(ctx context.Context, token string) (*model.TokenInfo, error) {
	info, err := l.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func encrypt(key, text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))
	return ciphertext, nil
}

func (l *Logic) CreateAuthToken(secret, str string) (res string, err error) {
	str = fmt.Sprintf("%s#%d", str, time.Now().UTC().UnixNano())
	b, err := encrypt([]byte(secret), []byte(str))
	return hex.EncodeToString(b), err
}
