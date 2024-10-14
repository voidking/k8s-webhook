# k8s-webhook
K8S Validating and Mutating Webhook Framework

Reference:
- [Go 语言设计与实现](https://draveness.me/golang/)
- [Gin 官方文档](https://gin-gonic.com/zh-cn/docs/introduction/)
- [viper](https://github.com/spf13/viper)

# Features
- When creating a pod, modify the pod definition: image, env, labels, etc.
- TODO: xxx 
- TODO: yyy

# Deploy to K8S
## Deploy all components
```bash
# ./scripts/uninstall-dev.sh
./scripts/install.sh
```

## Test Webhook Service
TODO

# Development
## Runtime
- go version 1.22.2
- kubernetes version 1.20.14 - 1.24.8

## Install dependencies
```bash
go mod tidy
```

## Unit Test
```bash
go test -run TestServeReturnsCorrectJson  k8s-webhook/internal/server
```

## Prepare K8S Resources
```bash
./scripts/install-dev.sh
```

## Run Webhook Service Locally
1. Configure environment
```bash
mkdir -p $HOME/cert && cp deployments/cert/* $HOME/cert
export HTTPS_CERT=$HOME/cert/cert.pem HTTPS_KEY=$HOME/cert/key.pem
```

2. Run Controller
```bash
go run main.go
```

## Test Webhook Service
1. curl http server
```bash
curl http://localhost:8080
curl -X POST http://localhost:8080 -H "Content-type: application/json" -d@test/data/create-ns-webhook.json
curl -X POST http://localhost:8080/pod/mutating -H "Content-type: application/json" -d@test/data/create-pod-webhook.json
```

2. curl https server
```bash
curl https://localhost:8443 --cacert $HOME/cert/ca.pem
curl -X POST https://localhost:8443 --cacert $HOME/cert/ca.pem -H "Content-type: application/json" -d@test/data/create-ns-webhook.json
curl -X POST https://localhost:8443/pod/mutating --cacert $HOME/cert/ca.pem -H "Content-type: application/json" -d@test/data/create-pod-webhook.json
```

## Build by Yourself
Build docker image
```bash
docker build -t k8s-webhook:latest .
```

PS: To accelerate build in China
```bash
docker build \
  --network=host \
  --build-arg HTTP_PROXY=http://192.168.56.1:7890 \
  --build-arg HTTPS_PROXY=http://192.168.56.1:7890 \
  -t k8s-webhook:latest .
```

# How to Prepare a TLS Certificate
k8s-webhook has been configured a self-signed TLS with ten years expiration.   
If you want to use your own certificate, you can follow the steps below.

1. Generate a certificate and private key
```bash
mkcert webhook.k8s-webhook.svc localhost
```

2. Download domain certificate and CA certificate, then rename
- ./webhook.k8s-webhook.svc+1.pem -> k8s-webhook/deployments/cert/cert.pem
- ./webhook.k8s-webhook.svc+1-key.pem -> k8s-webhook/deployments/cert/key.pem
- /root/.local/share/mkcert/rootCA.pem -> k8s-webhook/deployments/cert/ca.pem

