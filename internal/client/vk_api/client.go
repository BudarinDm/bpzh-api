package vk_api

import (
	"bpzh-api/internal/config"
	"bpzh-api/internal/model"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"strconv"
)

type Client struct {
	vk     *api.VK
	config *config.VkApiConfig
}

func NewClient(cfg *config.VkApiConfig) *Client {
	return &Client{
		config: cfg,
		vk:     api.NewVK(cfg.BotToken),
	}
}

func (c *Client) SendCode(code, userId int) (err error) {
	b := params.NewMessagesSendBuilder()

	b.Message(fmt.Sprintf("%d", code))
	b.RandomID(0)
	b.PeerID(userId)

	st, err := c.vk.MessagesSend(b.Params)
	fmt.Println(st)
	if err != nil {
		return err
	}
	return
}

func (c *Client) GetUsersInChat(chatId int64) (users []model.User, err error) {
	members, err := c.vk.MessagesGetConversationMembers(api.Params{"peer_id": chatId})
	if err != nil {
		return
	}

	var usersIds string
	for _, m := range members.Items {
		if m.MemberID < 0 {
			continue
		}
		usersIds += "," + strconv.Itoa(m.MemberID)
	}
	resp, err := c.vk.UsersGet(api.Params{"user_ids": usersIds, "fields": "domain"})
	if err != nil {
		return
	}

	for _, u := range resp {
		user := model.User{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			VkId:      int64(u.ID),
			VkDomain:  u.Domain,
		}
		users = append(users, user)
	}
	return
}
