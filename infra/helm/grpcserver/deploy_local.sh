#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --docker-username)
      dockerUsername="${2}"
      shift
      ;;
    --project)
      project="${2}"
      shift
      ;;
    --instance)
      instance="${2}"
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

# Language
if [[ $language == "" ]]; then
  echo -e "Language [--language] is not provided!\n"
  exit 1
fi

### Set variables

# redis
declare -A redis
redis["name"]="redis"
redis["namespace"]="ops"
redis["port"]=6379
redis["password"]="megasecret"

# otelcollectors
declare -A otelcollectors
otelcollectors["name"]="nrotelk8s"
otelcollectors["namespace"]="ops"
otelcollectors["endpoint"]="http://${otelcollectors[name]}-dep-rec-collector-headless.${otelcollectors[namespace]}.svc.cluster.local:4317"

# grpcserver
declare -A grpcserver
grpcserver["name"]="grpcserver"
grpcserver["imageName"]="${dockerUsername}/${project}-${grpcserver[name]}-${language}:latest"
grpcserver["namespace"]="${language}"
grpcserver["replicas"]=2
grpcserver["port"]=8080

###################
### Deploy Helm ###
###################

# grpcserver
helm upgrade ${grpcserver[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${grpcserver[namespace]} \
  --set imageName=${grpcserver[imageName]} \
  --set imagePullPolicy="Always" \
  --set language=${language} \
  --set name=${grpcserver[name]} \
  --set replicas=${grpcserver[replicas]} \
  --set port=${grpcserver[port]} \
  --set redis.server="${redis[name]}-master-0.${redis[name]}-headless.${redis[namespace]}.svc.cluster.local" \
  --set redis.port=${redis[port]} \
  --set redis.password="${redis[password]}" \
  --set otel.exporter="otlp" \
  --set otlp.endpoint="${otelcollectors[endpoint]}" \
  "./chart"
