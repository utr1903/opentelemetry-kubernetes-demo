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
    --k8s-version)
      k8sVersion="${2}"
      shift
      ;;
    --destroy)
      flagDestroy="true"
      shift
      ;;
    *)
      shift
      ;;
  esac
done

clusterName="kind${project}main${instance}"

if [[ $flagDestroy != "true" ]]; then
  kind create cluster \
    --name $clusterName \
    --config ../config/kind-config.yaml \
    --image=kindest/node:v${k8sVersion}
else
  kind delete cluster \
    --name $clusterName
fi
