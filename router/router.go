package router

import (
	"github.com/gin-contrib/cors"
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
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

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
