#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --github-actor)
      githubActor="${2}"
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

# redis
declare -A redis
redis["name"]="redis"
redis["namespace"]="ops"
redis["port"]=6379
redis["password"]="megasecret"

# mongo
declare -A mongo
mongo["name"]="mongo"
mongo["namespace"]="ops"
mongo["username"]="root"
mongo["password"]="verysecretpassword"
mongo["port"]=27017
mongo["database"]="otel"
mongo["table"]="${language}"

# otelcollectors
declare -A otelcollectors
otelcollectors["name"]="nrotelk8s"
otelcollectors["namespace"]="ops"
otelcollectors["endpoint"]="http://${otelcollectors[name]}-dep-rec-collector-headless.${otelcollectors[namespace]}.svc.cluster.local:4317"

# grpcserver
declare -A grpcserver
grpcserver["name"]="grpcserver"
grpcserver["imageName"]="ghcr.io/${githubActor}/${project}-${grpcserver[name]}-${language}:latest"
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
  --set redis.server="${redis[name]}-replicas.${redis[namespace]}.svc.cluster.local" \
  --set redis.port=${redis[port]} \
  --set redis.password="${redis[password]}" \
  --set mongo.server="${mongo[name]}.${mongo[namespace]}.svc.cluster.local" \
  --set mongo.username=${mongo[username]} \
  --set mongo.password=${mongo[password]} \
  --set mongo.port=${mongo[port]} \
  --set mongo.database=${mongo[database]} \
  --set mongo.table=${mongo[table]} \
  --set otel.exporter="otlp" \
  --set otlp.endpoint="${otelcollectors[endpoint]}" \
  "./infra/helm/${application}/chart"
