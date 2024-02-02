#!/bin/bash

### Set variables

# oteloperator
declare -A oteloperator
oteloperator["name"]="oteloperator"
oteloperator["namespace"]="ops"

###################
### Deploy Helm ###
###################

# Add helm repos
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update

# oteloperator
helm upgrade ${oteloperator[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${oteloperator[namespace]} \
  --version "0.40.0" \
  "open-telemetry/opentelemetry-operator"
