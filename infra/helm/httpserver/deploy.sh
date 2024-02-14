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

# mysql
declare -A mysql
mysql["name"]="mysql"
mysql["namespace"]="ops"
mysql["username"]="root"
mysql["password"]="verysecretpassword"
mysql["port"]=3306
mysql["database"]="otel"
mysql["table"]="${language}"

# otelcollectors
declare -A otelcollectors
otelcollectors["name"]="nrotelk8s"
otelcollectors["namespace"]="ops"
otelcollectors["endpoint"]="http://${otelcollectors[name]}-dep-rec-collector-headless.${otelcollectors[namespace]}.svc.cluster.local:4317"

# httpserver
declare -A httpserver
httpserver["name"]="httpserver"
httpserver["imageName"]="ghcr.io/${githubActor}/${project}-${httpserver[name]}-${language}:latest"
httpserver["namespace"]="${language}"
httpserver["replicas"]=2
httpserver["port"]=8080

###################
### Deploy Helm ###
###################

# httpserver
helm upgrade ${httpserver[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${httpserver[namespace]} \
  --set imageName=${httpserver[imageName]} \
  --set imagePullPolicy="Always" \
  --set language=${language} \
  --set name=${httpserver[name]} \
  --set replicas=${httpserver[replicas]} \
  --set port=${httpserver[port]} \
  --set redis.server="${redis[name]}-replicas.${redis[namespace]}.svc.cluster.local" \
  --set redis.port=${redis[port]} \
  --set redis.password="${redis[password]}" \
  --set mysql.server="${mysql[name]}.${mysql[namespace]}.svc.cluster.local" \
  --set mysql.username=${mysql[username]} \
  --set mysql.password=${mysql[password]} \
  --set mysql.port=${mysql[port]} \
  --set mysql.database=${mysql[database]} \
  --set mysql.table=${mysql[table]} \
  --set otel.exporter="otlp" \
  --set otlp.endpoint="${otelcollectors[endpoint]}" \
  "./infra/helm/${application}/chart"
