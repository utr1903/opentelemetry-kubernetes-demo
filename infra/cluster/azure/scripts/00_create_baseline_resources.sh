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
resourceGroupName="rg${project}base${instance}"
storageAccountName="st${project}base${instance}"
blobContainerName="${project}tfstates"

# Resource group
echo "Checking base resource group [${resourceGroupName}]..."
resourceGroup=$(az group show \
  --name $resourceGroupName \
  2> /dev/null)

if [[ $resourceGroup == "" ]]; then
  echo " -> Base resource group does not exist. Creating..."

  resourceGroup=$(az group create \
    --name $resourceGroupName \
    --location $location)

  echo -e " -> Base resource group is created successfully.\n"
else
  echo -e " -> Base resource group already exists.\n"
fi

# Storage account
echo "Checking base storage account [${storageAccountName}]..."
storageAccount=$(az storage account show \
    --resource-group $resourceGroupName \
    --name $storageAccountName \
  2> /dev/null)

if [[ $storageAccount == "" ]]; then
  echo " -> Base storage account does not exist. Creating..."

  storageAccount=$(az storage account create \
    --resource-group $resourceGroupName \
    --name $storageAccountName \
    --sku "Standard_LRS" \
    --allow-blob-public-access true \
    --encryption-services "blob")

  echo -e " -> Base storage account is created successfully.\n"
else
  echo -e " -> Base storage account already exists.\n"
fi

# Terraform blob container
echo "Checking Terraform blob container [${blobContainerName}]..."
terraformBlobContainer=$(az storage container show \
  --account-name $storageAccountName \
  --name $blobContainerName \
  2> /dev/null)

if [[ $terraformBlobContainer == "" ]]; then
  echo " -> Terraform blob container does not exist. Creating..."

  terraformBlobContainer=$(az storage container create \
    --account-name $storageAccountName \
    --name $blobContainerName \
    2> /dev/null)

  echo -e " -> Terraform blob container is created successfully.\n"
else
  echo -e " -> Terraform blob container already exists.\n"
fi
