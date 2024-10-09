package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/mattbaird/jsonpatch"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type PodAdmission struct {
}

func (*PodAdmission) HandleMutatingAdmission(c *gin.Context) {
	var admissionReview v1.AdmissionReview

	if err := c.ShouldBindJSON(&admissionReview); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pod := corev1.Pod{}
	if err := runtime.DecodeInto(Codecs.UniversalDeserializer(), []byte(admissionReview.Request.Object.Raw), &pod); err != nil {
		return
	}

	for key, value := range pod.Labels {
		logrus.Debugf("Key: %s, Value: %s", key, value)
	}

	// 修改pod定义
	for i := range pod.Spec.Containers {
		// 当前只处理第一个容器
		if i != 0 {
			break
		}
		container := &pod.Spec.Containers[i]
		container.Image = "newimage"
	}

	// 构造response
	resp := &v1.AdmissionResponse{
		UID:     admissionReview.Request.UID,
		Allowed: false,
	}

	originalPodBytes := admissionReview.Request.Object.Raw
	modifiedPodBytes, _ := json.Marshal(pod)
	patches, _ := jsonpatch.CreatePatch(originalPodBytes, modifiedPodBytes)
	patchesBytes, _ := json.Marshal(patches)
	patchType := v1.PatchTypeJSONPatch

	resp.PatchType = &patchType
	resp.Patch = patchesBytes
	resp.Allowed = true
	resp.Result = &metav1.Status{Status: "Success"}

	admissionReview.Response = resp
	c.JSON(200, admissionReview)
}

func (*PodAdmission) HandleValidatingAdmission(c *gin.Context) {
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
