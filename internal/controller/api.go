package controller

import (
	"bpzh-api/internal/config"
	"bpzh-api/internal/logic"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ErrorResponse типовой ответ с ошибкой
type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListData содержит список элементов ответа
type ListData struct {
	Items interface{} `json:"items"`
}

// ListRange описывает пейджер в ответе
type ListRange struct {
	Count  int64 `json:"count"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

// Response ответ со списком элементов
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// App основная структура для приложения
type App struct {
	router *gin.Engine
	config *config.Config
	logic  *logic.Logic
}

func (a *App) setV1Routes(router *gin.RouterGroup) {
	v1 := router.Group("/v1")

	v1.POST("/request_code", a.RequestCode)
	v1.POST("/login", a.CheckCode)

	v1.GET("/plug", func(c *gin.Context) {
		c.JSON(200, "plug f")
	})
}

func NewApp(config *config.Config, logic *logic.Logic) *App {
	if os.Getenv("DEBUG") == "true" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	return &App{
		router: gin.Default(),
		config: config,
		logic:  logic,
	}
}

func (a *App) StartServe() {
	a.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization",
			"ResponseType", "accept", "origin", "Cache-Control", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	a.setV1Routes(a.router.Group("/api"))

	server := &http.Server{
		Addr:           ":" + a.config.App.Port,
		Handler:        a.router,
		ReadTimeout:    time.Second * 15,
		WriteTimeout:   time.Second * 15,
		MaxHeaderBytes: 1 << 20,
	}
	//контекст для ожидания сигнала с ОС
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//Запускаем сервер
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Service running error")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Server stop start")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Server forced to shutdown")
	}
	log.Info().Msg("Server stopped")
}

// Error возвращает пользователю ошибку
func (a *App) Error(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{
		Status:  status,
		Message: message,
	})
}

type ErrorLimitResponse struct {
	Status          int    `json:"status"`
	Message         string `json:"message"`
	TillNextRequest int64  `json:"till_next_request"`
}

func (a *App) ErrorLimit(ctx *gin.Context, status int, message string, nextTry int64) {
	ctx.JSON(status, ErrorLimitResponse{
		Status:          status,
		Message:         message,
		TillNextRequest: nextTry,
	})
}

type CheckSessionError struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type CheckSessionResponse struct {
	Id   string `json:"id"`
	VkId int64  `json:"vk_id"`
}
