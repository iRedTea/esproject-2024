package handler

import (
	"esproject/docs"
	"esproject/pkg/repository"
	"esproject/pkg/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
	repo     *repository.Repository
}

func NewHandler(services *service.Service, repo *repository.Repository) *Handler {
	return &Handler{services: services, repo: repo}
}

//	@title			EasyStartup Project
//	@version		1.0
//	@description	Project service for EasyStartup.
//
//	@contact.name	API Support
//	@contact.email	red.tea.dev@gmail.com
//
//	@host		project.easystartup.su
//	@BasePath	/api
//

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	docs.SwaggerInfo.BasePath = "/"

	api := router.Group("/api", h.userIdentity)
	{
		api.DELETE("project/:id", h.delete)
		api.GET("project/:id", h.getProject)
		api.POST("project/:id", h.editProject)
		api.GET("project/all", h.all)
		api.GET("project/count", h.count)
		api.GET("project/my", h.my)
		api.PUT("project", h.create)
	}

	member := router.Group("/api/member", h.userIdentity)
	{
		member.GET("invite/:id", h.getInvite)
		member.GET("invite/my", h.myInvites)
		member.POST("invite/accept/:id", h.accept)
		member.POST("invite", h.invite)
		member.POST("kick", h.kick)

	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
