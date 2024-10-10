#!/bin/bash
# Delete MutatingWebhookConfiguration
kubectl delete -f ../deployments/manifest/mwc.yaml

# Delete Service
kubectl delete -f ../deployments/manifest/localservice.yaml

# Delete namespace
kubectl delete ns k8s-webhook