package logic

import (
	"bpzh-api/internal/model"
	"context"
)

func (l *Logic) GetUserByVkDomain(ctx context.Context, domain string) (*model.User, error) {
	return l.repo.GetUserByVkDomain(ctx, domain)
}
