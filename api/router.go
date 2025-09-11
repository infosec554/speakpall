package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "speakpall/api/docs"
	"speakpall/api/handler"
	"speakpall/pkg/logger"
	"speakpall/service"
)

// @title           Auth API
// @version         1.0
// @description     Authentication and Authorization API
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New(services service.IServiceManager, log logger.ILogger) *gin.Engine {
	h := handler.New(services, log)
	r := gin.New()

	// Global middlewares
	r.Use(gin.Recovery())

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	{
		auth.POST("/signup", h.SignUp)
		auth.POST("/login", h.Login)


	

	}

	

	return r
}
