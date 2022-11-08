package controller

import (
	"bpzh-api/internal/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"math/rand"
	"net/http"
	"time"
)

func (a *App) SetAuthRoutes(auth *gin.RouterGroup) {
	auth.POST("/request_code", a.RequestCode)
	auth.POST("/login", a.CheckCode)
	auth.GET("/check", a.CheckSession)
}

func (a *App) RequestCode(ctx *gin.Context) {
	// Валидация запроса
	var request model.RequestCode
	err := ctx.BindJSON(&request)
	if err != nil {
		a.Error(ctx, http.StatusBadRequest, "invalid_request")
		return
	}

	if request.Scope != "tg" && request.Scope != "vk" {
		a.Error(ctx, http.StatusBadRequest, "invalid_request")
		return
	}

	// Проверка предыдущих отправок
	lastSend, err := a.logic.GetLastSend(ctx, request)
	if err != nil {
		log.Info().Msgf("Ошибка проверки последней записи кода id=%s scope=%s : %s", request.Login, request.Scope, err.Error())
		a.Error(ctx, http.StatusInternalServerError, "err_check_last_code")
		return
	}
	fromLastSend := time.Now().Unix() - lastSend
	tillNextSend := 120 - fromLastSend

	if fromLastSend < 120 { // С момента предыдущей отправки прошло слишком мало времени
		log.Info().Msgf("Следующая попытка отправки кода на id %s через %d сек", request.Login, tillNextSend)
		a.ErrorLimit(ctx, http.StatusForbidden, "too_many_attempts", tillNextSend)
		return
	}

	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000

	user, err := a.logic.GetUserByVkDomain(ctx, request.Login)
	if err != nil {
		log.Info().Msgf("Ошибка запроса id=%s scope=%s : %s", request.Login, request.Scope, err.Error())
		a.Error(ctx, http.StatusInternalServerError, "err_get_user")
		return
	}

	log.Info().Msgf("Отправляем code для scope=%s id=%s", request.Scope, request.Login)
	//отправка кода через бота
	err = a.logic.SendCode(code, int(user.VkId))
	if err != nil {
		log.Info().Msgf("Ошибка отправки кода id=%s scope=%s : %s", request.Login, request.Scope, err.Error())
		a.Error(ctx, http.StatusInternalServerError, "err_send_code")
		return
	}

	// Сохранение информации об отправленном коде в нашей БД
	err = a.logic.UpdateCode(ctx, fmt.Sprintf("%s.%s", request.Scope, request.Login), int64(code))
	if err != nil {
		log.Error().Msgf("Ошибка сохранения кода для id=%s scope=%s : %s", request.Login, request.Scope, err.Error())
		a.Error(ctx, http.StatusInternalServerError, "save_code")
		return
	}

	resp := model.RequestCode{
		Login:           request.Login,
		TillNextRequest: 120,
	}

	ctx.JSON(http.StatusOK, resp)
}

func (a *App) CheckCode(ctx *gin.Context) {
	// Валидация запроса
	var request model.CheckCodeRequest
	err := ctx.BindJSON(&request)
	if err != nil {
		a.Error(ctx, http.StatusBadRequest, "invalid_request")
		return
	}

	if len(request.Login) == 0 {
		a.Error(ctx, http.StatusBadRequest, "invalid_token_length")
		return
	}

	// Ищем запись с code и достаем его

	code, tryCount, err := a.logic.GetCodeCheck(ctx, request.Scope, request.Login)
	if err != nil {
		return
	}

	if tryCount == 0 {
		log.Error().Msgf("Попыток подбора кода для domain=%s не осталось", request.Login)
		a.Error(ctx, http.StatusForbidden, "no_try_count")
		return
	}

	if request.Code != code {
		log.Error().Msgf("Неверный код для токена")
		a.Error(ctx, http.StatusForbidden, "incorrect_code")
		return
	}

	// Создаем новый токен для сессионной куки
	token, err := a.logic.CreateAuthToken(a.config.App.TokenSecret, request.Login)
	if err != nil {
		log.Error().Msgf("Ошибка формирования постоянного токена для domain %s: %s", request.Login, err)
		a.Error(ctx, http.StatusInternalServerError, "get_user_id_by_domain_failed")
		return
	}

	user, err := a.logic.GetUserByVkDomain(ctx, request.Login)
	if err != nil {
		log.Info().Msgf("Ошибка запроса id=%s scope=%s : %s", request.Login, request.Scope, err.Error())
		a.Error(ctx, http.StatusInternalServerError, "err_get_user")
		return
	}

	//пишем токен в базу
	err = a.logic.CreateToken(ctx, fmt.Sprintf("%s.%s", request.Scope, request.Login), token, user.DocId, user.VkId)
	if err != nil {
		log.Error().Msgf("Ошибка cоздания токена для domain %s: %s", request.Login, err)
		a.Error(ctx, http.StatusInternalServerError, "err_create_token")
		return
	}
	ctx.SetCookie(CookieName, token, 10080, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.DocId,
		"vk_id": user.VkId,
	})
}

func (a *App) CheckSession(ctx *gin.Context) {
	token, err := ctx.Cookie(CookieName)
	if err != nil {
		log.Error().Msgf("Токен не найден (cookie=%s): %s", CookieName, err)
		ctx.JSON(http.StatusOK, CheckSessionError{
			Status:  http.StatusBadRequest,
			Code:    ErrorCookieMissing,
			Message: "Cookie " + CookieName + " is missing",
		})
		return
	}

	tInfo, err := a.logic.GetSessionByToken(ctx, token)
	if err != nil {
		log.Error().Msgf("Сессия %s не найдена: %s", token, err)
		ctx.JSON(http.StatusOK, CheckSessionError{
			Status:  http.StatusForbidden,
			Code:    ErrorSessionNotFound,
			Message: "Cookie " + CookieName + " is missing",
		})
		return
	}

	ctx.SetCookie(CookieName, token, 10080, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"id":    tInfo.Id,
		"vk_id": tInfo.VkId,
	})
}
