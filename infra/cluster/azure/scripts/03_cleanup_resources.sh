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

### Set variables
baseResourceGroupName="rg${project}base${instance}"
servicePrincipalName="sp${project}${instance}"

# Service principal
echo "Checking service principal [${servicePrincipalName}]..."
subscriptionId=$(az account show | jq -r .id)

servicePrincipalAppId=$(az ad app list \
  --display-name $servicePrincipalName \
  2> /dev/null | jq -r .[0].appId)

if [[ $servicePrincipalAppId == "" ]]; then
  echo -e " -> Service principal does not exist.\n"
else
  echo -e " -> Service principal exists. Deleting..."
  
  az ad app delete \
    --id $servicePrincipalAppId
  
  echo -e " -> Service principal is deleted successfully.\n"
fi

# Resource group
echo "Checking base resource group [${baseResourceGroupName}]..."
resourceGroup=$(az group show \
  --name $baseResourceGroupName \
  2> /dev/null)

if [[ $resourceGroup == "" ]]; then
  echo -e " -> Resource group does not exist.\n"
else
  echo -e " -> Resource group exists. Deleting..."

  az group delete \
    --name $baseResourceGroupName \
    --yes

  echo -e " -> Resource group is deleted successfully.\n"
fi
