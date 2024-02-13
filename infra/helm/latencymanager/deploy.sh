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

# otelcollectors
declare -A otelcollectors
otelcollectors["name"]="nrotelk8s"
otelcollectors["namespace"]="ops"
otelcollectors["endpoint"]="http://${otelcollectors[name]}-dep-rec-collector-headless.${otelcollectors[namespace]}.svc.cluster.local:4317"

# latencymanager
declare -A latencymanager
latencymanager["name"]="latencymanager"
latencymanager["imageName"]="${dockerUsername}/${project}-${latencymanager[name]}-${language}:latest"
latencymanager["namespace"]="${language}"

###################
### Deploy Helm ###
###################

# latencymanager
helm upgrade ${latencymanager[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${latencymanager[namespace]} \
  --set imageName=${latencymanager[imageName]} \
  --set imagePullPolicy="Always" \
  --set language=${language} \
  --set name=${latencymanager[name]} \
  --set redis.server="${redis[name]}-master-0.${redis[name]}-headless.${redis[namespace]}.svc.cluster.local" \
  --set redis.port=${redis[port]} \
  --set redis.password="${redis[password]}" \
  --set otel.exporter="otlp" \
  --set otlp.endpoint="${otelcollectors[endpoint]}" \
  "./infra/helm/${application}/chart"
