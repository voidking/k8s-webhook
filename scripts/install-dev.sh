#!/bin/bash

# Run the script in project root path

# Create namespace
kubectl create ns k8s-webhook

# Create MutatingWebhookConfiguration
kubectl apply -f ./deployments/manifest/mwc.yaml

# Create Service
# vim manifest/localservice.yaml # change to local IP
kubectl apply -f ./deployments/manifest/localservice.yaml