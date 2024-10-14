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
	for i := range podTemplate.Spec.Containers {
		// 只处理第一个容器
		// if i != 0 {
		// 	break
		// }

		// pod模板中有多少个容器，就修改多少个容器
		if i >= len(pod.Spec.Containers) {
			break
		}
		containerTpl := &podTemplate.Spec.Containers[i]
		// 修改容器名和镜像
		container := &pod.Spec.Containers[i]
		container.Name = containerTpl.Name
		container.Image = containerTpl.Image
		// 修改env，无则新增，有则修改
		// logrus.Debugln("template env:", tplContainer.Env)
		// logrus.Debugln("env:", container.Env)
		for _, tplEnv := range containerTpl.Env {
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
		container.Command = containerTpl.Command
		container.Args = containerTpl.Args
		// 修改资源限制
		container.Resources = containerTpl.Resources

	}

	// 修改 labels ，无则新增，有则修改
	for tplKey, tplValue := range podTemplate.ObjectMeta.Labels {
		found := false
		for key := range pod.ObjectMeta.Labels {
			if tplKey == key {
				pod.ObjectMeta.Labels[key] = tplValue
				found = true
				break
			}
		}
		if !found {
			if pod.ObjectMeta.Labels == nil {
				pod.ObjectMeta.Labels = make(map[string]string)
			}
			pod.ObjectMeta.Labels[tplKey] = tplValue
		}
	}
	// logrus.Debugln("labels:", pod.ObjectMeta.Labels)

	// 修改 annotations ，无则新增，有则修改
	for tplKey, tplValue := range podTemplate.ObjectMeta.Annotations {
		found := false
		for key := range pod.ObjectMeta.Annotations {
			if tplKey == key {
				pod.ObjectMeta.Annotations[key] = tplValue
				found = true
				break
			}
		}
		if !found {
			if pod.ObjectMeta.Annotations == nil {
				pod.ObjectMeta.Annotations = make(map[string]string)
			}
			pod.ObjectMeta.Annotations[tplKey] = tplValue
		}
	}
	// logrus.Debugln("annotations:", pod.ObjectMeta.Annotations)

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
