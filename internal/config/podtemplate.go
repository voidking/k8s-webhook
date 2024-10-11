package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func GetPodTemplate() corev1.Pod {
	// 读取 podtemplate.yaml 文件
	podYamlPath := filepath.Join(GetConfigPath(), "podtemplate.yaml")
	bytes, err := os.ReadFile(podYamlPath)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// 解析 YAML 文件
	var pod corev1.Pod
	if err := yaml.Unmarshal(bytes, &pod); err != nil {
		logrus.Fatalf("Error parsing YAML: %v", err)
	}

	// 打印所需字段
	logrus.Infof("Pod Name: %s", pod.Name)
	for _, container := range pod.Spec.Containers {
		logrus.Infof("container.name: %s, container.image: %s", container.Name, container.Image)
	}
	return pod
}
