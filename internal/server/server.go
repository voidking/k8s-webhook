package server

import (
	"crypto/tls"
	"net/http"

	"k8s-webhook/internal/config"
	"k8s-webhook/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func GetAdmissionServerNoTLS(router *gin.Engine, Addr string) *http.Server {
	server := &http.Server{
		Handler: router,
		Addr:    Addr,
	}
	return server
}

func GetAdmissionValidationServer(router *gin.Engine, tlsCert, tlsKey, Addr string) *http.Server {
	sCert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
	server := GetAdmissionServerNoTLS(router, Addr)
	server.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{sCert},
	}
	if err != nil {
		logrus.Error(err)
	}
	return server
}

var (
	g errgroup.Group
)

func RunServer() {
	conf := config.GetConfig()
	logrus.Infoln(conf)

	// gin.SetMode(gin.ReleaseMode)
	serverRouter, _ := router.SetRouter()

	// Run server without tls
	server_http := GetAdmissionServerNoTLS(serverRouter, conf.HTTP.Addr)
	// server_http.ListenAndServe()

	g.Go(func() error {
		return server_http.ListenAndServe()
	})

	if conf.HTTPS.Enable {
		// Run server with tls
		server_https := GetAdmissionValidationServer(serverRouter, conf.HTTPS.Cert, conf.HTTPS.Key, conf.HTTPS.Addr)
		// server_https.ListenAndServeTLS("", "")
		g.Go(func() error {
			return server_https.ListenAndServeTLS("", "")
		})
	}

	if err := g.Wait(); err != nil {
		logrus.Error(err)
	}
}
