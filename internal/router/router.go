package router

import (
	"k8s-webhook/internal/api"

	"github.com/gin-gonic/gin"
)

func SetRouter() (*gin.Engine, error) {
	router := gin.Default()
	nsac := api.NamespaceAdmission{}
	podac := api.PodAdmission{}

	router.Any("/", nsac.HandleAdmission)
	podRouter := router.Group("/pod")
	podRouter.POST("/mutating", podac.HandleMutatingAdmission)
	podRouter.POST("/validating", podac.HandleValidatingAdmission)
	return router, nil
}
