package v1

import (
	"testovoye/internal/service"

	"github.com/gin-gonic/gin"
)

type HandlerV1 struct {
	services *service.Service
}

func NewHandlerV1(services *service.Service) *HandlerV1 {
	return &HandlerV1{
		services: services,
	}
}

func (h *HandlerV1) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initPersonRoutes(v1)
	}
}
