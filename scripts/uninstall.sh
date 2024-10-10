#!/bin/bash

# Delete MutatingWebhookConfiguration
kubectl delete -f ../deployments/manifest/mwc.yaml

# Delete Deployment
kubectl delete -f ../deployments/manifest/deployment.yaml
kubectl delete -f ../deployments/manifest/configmap/webhook-env-cm.yaml

# Delete Service
kubectl delete -f ../deployments/manifest/service.yaml

# Delete namespace
kubectl delete ns k8s-webhook

# Hint
echo -e "\n\n"
echo "Uninstall completed!"