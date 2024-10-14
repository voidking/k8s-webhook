package api

import (
	"encoding/json"
	"k8s-webhook/internal/config"

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

	// 验证pod label
	// for key, value := range pod.Labels {
	// 	logrus.Debugf("Key: %s, Value: %s", key, value)
	// }

	// 读取pod模板
	podTemplate := config.GetPodTemplate()

	// 修改pod定义
	for i := range pod.Spec.Containers {
		// 当前只处理第一个容器
		if i != 0 {
			break
		}
		tplContainer := &podTemplate.Spec.Containers[i]
		// 修改镜像
		container := &pod.Spec.Containers[i]
		container.Image = tplContainer.Image
		// 修改env，无则新增，有则修改
		// logrus.Debugln("template env:", tplContainer.Env)
		// logrus.Debugln("env:", container.Env)
		for _, tplEnv := range tplContainer.Env {
			found := false
			for j := range container.Env {
				env := &container.Env[j]
				if tplEnv.Name == env.Name {
					env.Value = tplEnv.Value
					found = true
					break
				}
			}
			if !found {
				container.Env = append(container.Env, tplEnv)
			}
		}
		// 修改启动命令，包括 command 和 args
		container.Command = tplContainer.Command
		container.Args = tplContainer.Args
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
