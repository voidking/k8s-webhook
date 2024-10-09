package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NamespaceAdmission struct {
}

func (*NamespaceAdmission) HandleAdmission(c *gin.Context) {
	var review v1.AdmissionReview
	c.ShouldBindJSON(&review)
	logrus.Debugln(review.Request)
	review.Response = &v1.AdmissionResponse{
		Allowed: true,
		Result: &metav1.Status{
			Message: "Welcome aboard!",
		},
	}
	c.JSON(200, review)
}
