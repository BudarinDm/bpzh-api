package controller

//func (a *App) checkAuth(ctx *gin.Context) {
//	cookieNameWb := auth.WBCookieName
//	cookieNameCrm := auth.CRMCookieName
//	var scope string
//
//	sessionToken, err := ctx.Cookie(cookieNameWb)
//	if err != nil {
//		sessionToken, err = ctx.Cookie(cookieNameCrm)
//		if err != nil {
//			log.Error().Msg(fmt.Sprintf("Сессионная кука %s не найдена: %s", cookieNameWb, err))
//			a.Error(ctx, http.StatusUnauthorized, ErrorAuthNoSessionCookie, "Unauthorized")
//			ctx.Abort()
//			return
//		}
//		scope = "crm"
//	}
//	wbuserId, err := a.logic.CheckAuth(ctx, sessionToken, scope)
//	if err != nil {
//		log.Error().Msg(fmt.Sprintf("Ошибка проверки сессии: %s", err))
//		a.Error(ctx, http.StatusUnauthorized, ErrorAuthServiceError, "Unauthorized")
//		ctx.Abort()
//		return
//	}
//
//	log.Info().Msg(fmt.Sprintf("Аутентифицирован пользователь wbuserid[%d]", wbuserId))
//
//	ctx.Set("wbuser_id", wbuserId)
//	ctx.Next()
//}
