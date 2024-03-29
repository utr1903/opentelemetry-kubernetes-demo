#!/bin/bash

### Set variables
declare -A redis
redis["name"]="redis"
redis["namespace"]="ops"
redis["password"]="megasecret"

###################
### Deploy Helm ###
###################

# Add helm repos
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# kafka
helm upgrade ${redis[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${redis[namespace]} \
  --set auth.password=${redis[password]} \
  --version "18.12.1" \
  "bitnami/redis"
