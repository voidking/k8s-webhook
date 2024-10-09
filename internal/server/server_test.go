package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s-webhook/internal/api"
)

var (
	AdmissionRequestNS = v1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind: "AdmissionReview",
		},
		Request: &v1.AdmissionRequest{
			UID: "e911857d-c318-11e8-bbad-025000000001",
			Kind: metav1.GroupVersionKind{
				Kind: "Namespace",
			},
			Operation: "CREATE",
			Object: runtime.RawExtension{
				Raw: []byte(`{"metadata": {
        						"name": "test",
        						"uid": "e911857d-c318-11e8-bbad-025000000001",
						        "creationTimestamp": "2018-09-28T12:20:39Z"
      						}}`),
			},
		},
	}
)

func decodeResponse(body io.ReadCloser) *v1.AdmissionReview {
	bodyBytes, _ := io.ReadAll(body)
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	review := &v1.AdmissionReview{}
	api.Codecs.UniversalDeserializer().Decode(bodyBytes, nil, review)
	return review
}

func encodeRequest(review *v1.AdmissionReview) []byte {
	ret, err := json.Marshal(review)
	if err != nil {
		logrus.Errorln(err)
	}
	return ret
}

func TestServeReturnsCorrectJson(t *testing.T) {
	router := gin.Default()
	nsac := api.NamespaceAdmission{}
	router.Any("/", nsac.HandleAdmission)
	server := httptest.NewServer(GetAdmissionServerNoTLS(router, ":8080").Handler)
	requestString := string(encodeRequest(&AdmissionRequestNS))
	myr := strings.NewReader(requestString)
	r, _ := http.Post(server.URL, "application/json", myr)
	review := decodeResponse(r.Body)

	if review.Request.UID != AdmissionRequestNS.Request.UID {
		t.Error("Request and response UID don't match")
	}
}
