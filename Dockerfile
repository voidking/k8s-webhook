FROM golang:1.22.2-alpine as builder

RUN apk update && apk add git && apk add ca-certificates

WORKDIR /k8s-webhook

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/k8s-webhook

# Runtime image
# FROM scratch AS base
FROM debian:bookworm-slim AS base
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /k8s-webhook/deployments/cert/* /etc/k8s-webhook/cert/
COPY --from=builder /k8s-webhook/configs/* /etc/k8s-webhook/configs/
COPY --from=builder /go/bin/k8s-webhook /bin/k8s-webhook

ENV CONFIG_PATH=/etc/k8s-webhook/configs \
    HTTP_ADDR=0.0.0.0:8080 \
    HTTPS_ADDR=0.0.0.0:8443 \
    HTTPS_ENABLE=true \
    HTTPS_CERT=/etc/k8s-webhook/cert/cert.pem \
    HTTPS_KEY=/etc/k8s-webhook/cert/key.pem
EXPOSE 8080 8443
ENTRYPOINT ["/bin/k8s-webhook"]