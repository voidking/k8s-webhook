package config

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// go test -v -run TestGetPodTemplate k8s-webhook/internal/config
func TestGetPodTemplate(t *testing.T) {
	pod := GetPodTemplate()
	logrus.Info(pod)
}
