package handler

import (
	"context"
	_ "effectiveMobile/docs"
	"effectiveMobile/internal/logger"
	"effectiveMobile/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes(ctx context.Context) http.Handler {
	log := logger.LoggerFromContext(ctx)
	log.Debugw("initializing routes")
	router := gin.Default()
	router.Use(logger.LoggerMiddleware(log))
	log.Debugw("swagger is initialized")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api")
	{
		api.GET("/getsongs", h.GetSongs)
		api.GET("/gettext/:id", h.GetTextSong)
		api.DELETE("deletesong/:id", h.DeleteSong)
		api.PATCH("updatesong/:id", h.UpdateSong)
		api.POST("insertsong/", h.InsertSong)
	}
	log.Debugw("routes are initialized")
	return router
}
