package router

import (
	"k8s-webhook/internal/api"

	"github.com/gin-gonic/gin"
)

func SetRouter() (*gin.Engine, error) {
	// router := gin.Default()
	router := gin.New()
	// router.Use(gin.Logger())
	// router.Use(CustomLogger())
	router.Use(LoggerNoHealthz())

	nsac := api.NamespaceAdmission{}
	podac := api.PodAdmission{}

	// 根路由
	router.Any("/", nsac.HandleAdmission)
	// 探活路由
	router.GET("/healthz", api.Healthz)
	// pod相关路由
	podRouter := router.Group("/pod")
	podRouter.POST("/mutating", podac.HandleMutatingAdmission)
	podRouter.POST("/validating", podac.HandleValidatingAdmission)

	return router, nil
}
