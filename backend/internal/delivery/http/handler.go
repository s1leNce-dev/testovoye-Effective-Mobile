package http

import (
	v1 "testovoye/internal/delivery/http/v1"
	"testovoye/internal/service"

	_ "testovoye/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HandlerAPI struct {
	services *service.Service
}

func NewHandlerAPI(services *service.Service) *HandlerAPI {
	return &HandlerAPI{
		services: services,
	}
}

func (h *HandlerAPI) Init() *gin.Engine {
	router := gin.Default()

	router.Use(
		corsMiddleware,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.initAPI(router)

	return router
}

func (h *HandlerAPI) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandlerV1(h.services)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
