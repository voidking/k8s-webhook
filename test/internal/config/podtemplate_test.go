package test

import (
	"k8s-webhook/internal/config"
	"testing"

	"github.com/sirupsen/logrus"
)

// go test -v test/internal/config/podtemplate_test.go -run TestGetPodTemplate
func TestGetPodTemplate(t *testing.T) {
	pod := config.GetPodTemplate()
	logrus.Info(pod)
}
