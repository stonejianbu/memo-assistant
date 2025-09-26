package router

import (
	"github.com/gin-gonic/gin"
	"github.com/stonejianbu/memo-assistant/config"
	"github.com/stonejianbu/memo-assistant/dao"
	"github.com/stonejianbu/memo-assistant/handler"
	"github.com/stonejianbu/memo-assistant/service"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
)

// SetupRouter server router
func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	server := gin.New()

	// init middleware
	server.Use(gin.Recovery())

	// init dao
	dao.InitWeaviate(weaviate.Config{
		Host:   config.Cfg.Weaviate.Host,
		Scheme: config.Cfg.Weaviate.Schema,
	}, config.Cfg.Weaviate.Class)

	// init services
	service.Init()

	// init route
	api := server.Group("/api/v1")
	{
		api.POST("/generate", handler.RetrieveText)
		api.POST("/train", handler.TrainText)
	}
	return server
}
