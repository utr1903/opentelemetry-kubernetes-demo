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

# kafka
declare -A kafka
kafka["name"]="kafka"
kafka["namespace"]="ops"
kafka["topic"]="${language}"

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

# kafkaconsumer
declare -A kafkaconsumer
kafkaconsumer["name"]="kafkaconsumer"
kafkaconsumer["imageName"]="ghcr.io/${githubActor}/${project}-${kafkaconsumer[name]}-${language}:latest"
kafkaconsumer["namespace"]="${language}"
kafkaconsumer["replicas"]=2
kafkaconsumer["port"]=8080

###################
### Deploy Helm ###
###################

# kafkaconsumer
helm upgrade ${kafkaconsumer[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${kafkaconsumer[namespace]} \
  --set dockerhubName=$DOCKERHUB_NAME \
  --set imageName=${kafkaconsumer[imageName]} \
  --set imagePullPolicy="Always" \
  --set language=${language} \
  --set name=${kafkaconsumer[name]} \
  --set replicas=${kafkaconsumer[replicas]} \
  --set port=${kafkaconsumer[port]} \
  --set kafka.address="${kafka[name]}.${kafka[namespace]}.svc.cluster.local:9092" \
  --set kafka.topic=${kafka[topic]} \
  --set kafka.groupId=${kafkaconsumer[name]} \
  --set mysql.server="${mysql[name]}.${mysql[namespace]}.svc.cluster.local" \
  --set mysql.username=${mysql[username]} \
  --set mysql.password=${mysql[password]} \
  --set mysql.port=${mysql[port]} \
  --set mysql.database=${mysql[database]} \
  --set mysql.table=${mysql[table]} \
  --set otel.exporter="otlp" \
  --set otlp.endpoint="${otelcollectors[endpoint]}" \
  "./infra/helm/${application}/chart"
