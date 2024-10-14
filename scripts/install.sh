#!/bin/bash
# Create namespace
kubectl create ns k8s-webhook

# Create MutatingWebhookConfiguration
kubectl apply -f ../deployments/manifest/mwc.yaml

# Create Deployment
kubectl apply -f ../deployments/manifest/configmap/webhook-env.yaml
kubectl apply -f ../deployments/manifest/configmap/webhook-config.yaml
kubectl apply -f ../deployments/manifest/deployment.yaml

# Create Service
kubectl apply -f ../deployments/manifest/service.yaml