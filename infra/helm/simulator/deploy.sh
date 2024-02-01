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
    --application)
      application="${2}"
      shift
      ;;
    --language)
      language="${2}"
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

# Application
if [[ $application == "" ]]; then
  echo -e "Application [--application] is not provided!\n"
  exit 1
fi

# Language
if [[ $language == "" ]]; then
  echo -e "Language [--language] is not provided!\n"
  exit 1
fi

### Set variables

# kafka
declare -A kafka
kafka["name"]="kafka"
kafka["namespace"]="ops"
kafka["topic"]="${language}"

# otelcollectors
declare -A otelcollectors
otelcollectors["name"]="nrotelk8s"
otelcollectors["namespace"]="ops"
otelcollectors["endpoint"]="http://${otelcollectors[name]}-dep-rec-collector-headless.${otelcollectors[namespace]}.svc.cluster.local:4317"

# httpserver
declare -A httpserver
httpserver["name"]="httpserver"
httpserver["namespace"]="${language}"
httpserver["port"]=8080

# simulator
declare -A simulator
simulator["name"]="simulator-${language}"
simulator["imageName"]="ghcr.io/utr1903/${project}-${simulator[name]}-${language}:latest"
simulator["namespace"]="${language}"
simulator["replicas"]=3
simulator["port"]=8080
simulator["httpInterval"]=2000
simulator["kafkaInterval"]=1000

###################
### Deploy Helm ###
###################

# simulator
helm upgrade ${simulator[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${simulator[namespace]} \
  --set imageName=${simulator[imageName]} \
  --set imagePullPolicy="Always" \
  --set name=${simulator[name]} \
  --set replicas=${simulator[replicas]} \
  --set port=${simulator[port]} \
  --set httpserver.requestInterval=${simulator[httpInterval]} \
  --set httpserver.endpoint="${httpserver[name]}.${httpserver[namespace]}.svc.cluster.local" \
  --set httpserver.port="${httpserver[port]}" \
  --set kafka.address="${kafka[name]}.${kafka[namespace]}.svc.cluster.local:9092" \
  --set kafka.topic=${kafka[topic]} \
  --set kafka.requestInterval=${simulator[kafkaInterval]} \
  --set otel.exporter="otlp" \
  --set otlp.endpoint="${otelcollectors[endpoint]}" \
  "./infra/helm/${application}/chart"
