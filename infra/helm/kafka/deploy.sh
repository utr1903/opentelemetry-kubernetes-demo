#!/bin/bash

### Set variables
declare -A kafka
kafka["name"]="kafka"
kafka["namespace"]="opsteam"

###################
### Deploy Helm ###
###################

# Add helm repos
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# kafka
helm upgrade ${kafka[name]} \
  --install \
  --wait \
  --debug \
  --create-namespace \
  --namespace=${kafka[namespace]} \
  --set listeners.client.protocol=PLAINTEXT \
  --set provisioning.enabled=true \
  --set provisioning.topics[0].name="dotnet" \
  --set provisioning.topics[0].partitions=3 \
  --set provisioning.topics[1].name="golang" \
  --set provisioning.topics[1].partitions=3 \
  --set provisioning.topics[2].name="java" \
  --set provisioning.topics[2].partitions=3 \
  --set provisioning.topics[3].name="javascript" \
  --set provisioning.topics[4].partitions=3 \
  --set provisioning.topics[5].name="python" \
  --set provisioning.topics[5].partitions=3 \
  --version "26.6.2" \
  "bitnami/kafka"
