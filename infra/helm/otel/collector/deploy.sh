#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --project)
      project="${2}"
      shift
      ;;
    --instance)
      instance="${2}"
      shift
      ;;
    --cluster-type)
      clusterType="${2}"
      shift
      ;;
    --newrelic-otlp-endpoint)
      newrelicOtlpEndpoint="${2}"
      shift
      ;;
    --newrelic-opsteam-license-key)
      newrelicOpsteamLicenseKey="${2}"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

### Check input

# Project
if [[ $project == "" ]]; then
  echo -e "Project [--project] is not provided!\n"
  exit 1
fi

# Instance
if [[ $instance == "" ]]; then
  echo -e "Instance [--instance] is not provided!\n"
  exit 1
fi

# Cluster type
if [[ $clusterType == "" ]]; then
  echo "Cluster type [--cluster-type] is not given."
  exit 1
else
  if [[ $clusterType != "aks" ]]; then
    echo "Given cluster type [--cluster-type] is not supported. Supported values are: aks."
    exit 1
  fi
  clusterName="${clusterType}${project}${instance}"
fi

# New Relic OTLP endpoint
if [[ $newrelicOtlpEndpoint == "" ]]; then
  echo -e "New Relic OTLP endpoint [--newrelic-otlp-endpoint] is not provided.\n"
  exit 1
else
  if [[ $newrelicOtlpEndpoint != "otlp.nr-data.net:4317" && $newrelicOtlpEndpoint != "otlp.eu01.nr-data.net:4317" ]]; then
    echo "Given New Relic OTLP endpoint [--newrelic-otlp-endpoint] is not supported. Supported values are: US -> otlp.nr-data.net:4317, EU -> otlp.eu01.nr-data.net:4317."
    exit 1
  fi
fi

# New Relic OTLP endpoint
if [[ $newrelicOpsteamLicenseKey == "" ]]; then      
  echo -e "New Relic opsteam license key [--newrelic-opsteam-license-key] is not provided!\n"
  exit 1
fi

### Set variables

# otelcollectors
declare -A otelcollectors
otelcollectors["name"]="nrotelk8s"
otelcollectors["namespace"]="monitoring"

###################
### Deploy Helm ###
###################

# Add helm repos
helm repo add newrelic-experimental https://newrelic-experimental.github.io/monitoring-kubernetes-with-opentelemetry/charts
helm repo update

# otelcollector
helm upgrade ${otelcollectors[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace ${otelcollectors[namespace]} \
  --set clusterName=$clusterName \
  --set global.newrelic.enabled=true \
  --set global.newrelic.endpoint=$newrelicOtlpEndpoint \
  --set global.newrelic.teams.opsteam.licenseKey.value=$newrelicOpsteamLicenseKey \
  --version "0.6.0" \
  "newrelic-experimental/nrotelk8s"
