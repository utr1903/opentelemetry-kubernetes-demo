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
    --cluster-type)
      clusterType="${2}"
      shift
      ;;
    --newrelic-opsteam-account-id)
      newrelicOpsteamAccountId="${2}"
      shift
      ;;
    --newrelic-user-api-key)
      newrelicUserApiKey="${2}"
      shift
      ;;
    --newrelic-region)
      newrelicRegion="${2}"
      shift
      ;;
    --dry-run)
      flagDryRun="${2}"
      shift
      ;;
    --destroy)
      flagDestroy="${2}"
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

# New Relic OPS team account ID
if [[ $newrelicOpsteamAccountId == "" ]]; then
  echo -e "New Relic OPS team account ID [--newrelic-opsteam-account-id] is not provided!\n"
  exit 1
fi

# New Relic region
if [[ $newrelicRegion == "" ]]; then
  echo -e "New Relic region [--newrelic-region] is not provided.\n"
  exit 1
else
  if [[ $newrelicRegion != "us" && $newrelicRegion != "eu" ]]; then
    echo "Given New Relic region [--newrelic-region] is not supported. Supported values are: us & eu."
    exit 1
  fi
fi

# Cluster name
if [[ $clusterType == "" ]]; then
  echo "Cluster type [--cluster-type] is not given."
  exit 1
else
  if [[ $clusterType != "aks" ]]; then
    echo "Given cluster type [--cluster-type] is not supported. Supported values are: aks."
    exit 1
  fi
  clusterName="${clusterType}${project}${instance}"
fi

### Set variables
resourceGroupName="rg${project}base${instance}"
storageAccountName="st${project}base${instance}"
blobContainerName="${project}tfstates"

if [[ $flagDestroy != "true" ]]; then

  # Initialize Terraform
  terraform -chdir=./monitoring/newrelic/infra/terraform init \
    -backend-config="resource_group_name=${resourceGroupName}" \
    -backend-config="storage_account_name=${storageAccountName}" \
    -backend-config="container_name=${blobContainerName}" \
    -backend-config="key=monitoring-newrelic-infra"

  # Plan Terraform
  terraform -chdir=./monitoring/newrelic/infra/terraform plan \
    -var NEW_RELIC_ACCOUNT_ID=$newrelicOpsteamAccountId \
    -var NEW_RELIC_API_KEY=$newrelicUserApiKey \
    -var NEW_RELIC_REGION=$newrelicRegion \
    -var cluster_name=$clusterName \
    -out "./tfplan"

  # Apply Terraform
  if [[ $flagDryRun != "true" ]]; then
    terraform -chdir=./monitoring/newrelic/infra/terraform apply \
      -auto-approve \
      tfplan
  fi
else

  # Destroy Terraform
  terraform -chdir=./monitoring/newrelic/infra/terraform destroy \
    -auto-approve \
    -var NEW_RELIC_ACCOUNT_ID=$newrelicOpsteamAccountId \
    -var NEW_RELIC_API_KEY=$newrelicUserApiKey \
    -var NEW_RELIC_REGION=$newrelicRegion \
    -var cluster_name=$clusterName
fi
