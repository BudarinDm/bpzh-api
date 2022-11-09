package logic

import (
	"bpzh-api/internal/model"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (l *Logic) GetUserByVkDomain(ctx context.Context, domain string) (*model.User, error) {
	return l.repo.GetUserByVkDomain(ctx, domain)
}

func (l *Logic) UpdateUsersInFSVK(ctx context.Context, chatId int64, botId string) (err error) {
	var users []model.User
	if botId == "bpzh-vk-bot" {
		users, err = l.vkApiBpzhClient.GetUsersInChat(chatId)
		if err != nil {
			return err
		}
	}

	for _, u := range users {
		var domainUser *model.User
		domainUser, err = l.repo.GetUserByVkDomain(ctx, u.VkDomain)
		if err != nil {
			return
		}
		if domainUser.DocId == "" {
			domainUser, err = l.repo.GetUserByVkId(ctx, u.VkId)
			if err != nil && status.Code(err) != codes.NotFound {
				return
			}
			if domainUser.DocId == "" {
				err = l.repo.CreateUserOnVk(ctx, u)
				if err != nil {
					return err
				}
				continue
			}
			//обновляем юзера если смогли его найти
			err = l.repo.UpdateUserOnVk(ctx, domainUser.DocId, u)
			if err != nil {
				return err
			}
			continue
		}
		err = l.repo.UpdateUserOnVk(ctx, domainUser.DocId, u)
		if err != nil {
			return err
		}
	}

	return
}
