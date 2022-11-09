package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func (a *App) SetUserRoutes(g *gin.RouterGroup) {
	g.PUT("/vk-chat-update", a.UpdateUsersInFSVK)
}

func (a *App) UpdateUsersInFSVK(c *gin.Context) {
	var err error
	chatIdRow := c.Query("chat_id")
	chatId, err := strconv.ParseInt(chatIdRow, 10, 64)
	if err != nil {
		log.Error().Err(err).Msg("Invalid parse chat_id")
		a.Error(c, http.StatusBadRequest, "invalid_parse_chat_id")
		return
	}

	err = a.logic.UpdateUsersInFSVK(c, chatId, "bpzh-vk-bot")
	if err != nil {
		log.Error().Err(err).Msg("Error updating users")
		a.Error(c, http.StatusInternalServerError, "error_updating_users")
		return
	}

	c.JSON(http.StatusOK,
		Response{
			Status: http.StatusOK,
		},
	)

}
