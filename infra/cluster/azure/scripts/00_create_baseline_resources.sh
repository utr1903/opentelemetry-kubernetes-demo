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
baseStorageAccountName="st${project}base${instance}"
baseBlobContainerName="${project}tfstates"

# Resource group
echo "Checking base resource group [${baseResourceGroupName}]..."
baseResourceGroup=$(az group show \
  --name $baseResourceGroupName \
  2> /dev/null)

if [[ $baseResourceGroup == "" ]]; then
  echo " -> Base resource group does not exist. Creating..."

  baseResourceGroup=$(az group create \
    --name $baseResourceGroupName \
    --location $location)

  echo -e " -> Base resource group is created successfully.\n"
else
  echo -e " -> Base resource group already exists.\n"
fi

# Storage account
echo "Checking base storage account [${baseStorageAccountName}]..."
baseStorageAccount=$(az storage account show \
    --resource-group $resourceGroupName \
    --name $baseStorageAccountName \
  2> /dev/null)

if [[ $baseStorageAccount == "" ]]; then
  echo " -> Base storage account does not exist. Creating..."

  baseStorageAccount=$(az storage account create \
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
echo "Checking Terraform blob container [${baseBlobContainerName}]..."
terraformBlobContainer=$(az storage container show \
  --account-name $baseStorageAccountName \
  --name $baseBlobContainerName \
  2> /dev/null)

if [[ $terraformBlobContainer == "" ]]; then
  echo " -> Terraform blob container does not exist. Creating..."

  terraformBlobContainer=$(az storage container create \
    --account-name $baseStorageAccountName \
    --name $baseBlobContainerName \
    2> /dev/null)

  echo -e " -> Terraform blob container is created successfully.\n"
else
  echo -e " -> Terraform blob container already exists.\n"
fi
