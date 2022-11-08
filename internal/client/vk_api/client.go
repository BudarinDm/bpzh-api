package vk_api

import (
	"bpzh-api/internal/config"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
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
