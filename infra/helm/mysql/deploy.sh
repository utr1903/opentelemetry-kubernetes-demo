#!/bin/bash

### Set variables

# mysql
declare -A mysql
mysql["name"]="mysql"
mysql["namespace"]="ops"
mysql["password"]="verysecretpassword"
mysql["database"]="otel"

###################
### Deploy Helm ###
###################

# Add helm repos
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# mysql
helm upgrade ${mysql[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${mysql[namespace]} \
  --set auth.rootPassword=${mysql[password]} \
  --set auth.database=${mysql[database]} \
  --version "9.15.0" \
  "bitnami/mysql"
