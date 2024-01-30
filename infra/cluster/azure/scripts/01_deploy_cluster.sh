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
    --location)
      location="${2}"
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
    --dry-run)
      flagDryRun="true"
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

# Location
if [[ $location == "" ]]; then
  location="westeurope"
  echo -e "Location [--location] is not provided. Using default location ${location}.\n"
fi

# K8s version
if [[ $k8sVersion == "" ]]; then
  k8sVersion="1.28.0"
  echo -e "K8s version [--k8s-version] is not provided. Using default version ${k8sVersion}.\n"
fi

### Set variables

# Base
baseResourceGroupName="rg${project}base${instance}"
baseKeyVaultName="kv${project}base${instance}"
baseStorageAccountName="st${project}base${instance}"
baseBlobContainerName="${project}tfstates"

# Main
mainResourceGroupName="rg${project}main${instance}"
mainAksResourceName="aks${project}main${instance}"
mainAksNodepoolResourceGroupName="rgaks${project}main${instance}"

### Perform Terraform deployment
azureAccount=$(az account show)
tenantId=$(echo $azureAccount | jq -r .tenantId)
subscriptionId=$(echo $azureAccount | jq -r .id)

if [[ $flagDestroy != "true" ]]; then

  # Initialize Terraform
  terraform -chdir=../terraform init \
    -upgrade \
    -backend-config="tenant_id=${tenantId}" \
    -backend-config="subscription_id=${subscriptionId}" \
    -backend-config="resource_group_name=${baseResourceGroupName}" \
    -backend-config="storage_account_name=${baseStorageAccountName}" \
    -backend-config="container_name=${baseBlobContainerName}" \
    -backend-config="key=cluster.tfstate"

  # Plan Terraform
  terraform -chdir=../terraform plan \
    -var aks_resource_group_name=$mainResourceGroupName \
    -var aks_resource_name=$mainAksResourceName \
    -var aks_nodepool_resource_name=$mainAksNodepoolResourceGroupName \
    -var aks_version=$k8sVersion \
    -var location=$location \
    -out "./tfplan"

    if [[ $flagDryRun != "true" ]]; then
    
      # Apply Terraform
      terraform -chdir=../terraform apply tfplan

      # Get AKS credentials
      az aks get-credentials \
        --resource-group $mainResourceGroupName \
        --name $mainAksResourceName \
        --overwrite-existing
    fi
else

  # Destroy resources
  terraform -chdir=../terraform destroy \
    -var aks_resource_group_name=$mainResourceGroupName \
    -var aks_resource_name=$mainAksResourceName \
    -var aks_nodepool_resource_name=$mainAksNodepoolResourceGroupName \
    -var aks_version=$k8sVersion \
    -var location=$location
fi
