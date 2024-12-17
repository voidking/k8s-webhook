package config

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// go test -v test/internal/config/podtemplate_test.go -run TestGetPodTemplate
func TestGetPodTemplate(t *testing.T) {
	pod := GetPodTemplate()
	logrus.Info(pod)
}
