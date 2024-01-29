#!/bin/bash

# Get commandline arguments
while (( "$#" )); do
  case "$1" in
    --owner)
      owner="${2}"
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

# Owner
if [[ $owner == "" ]]; then
  echo -e "Owner [--owner] is not provided!\n"
  exit 1
fi

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
resourceGroupName="rg${owner}${project}${instance}"
keyVaultName="kv${owner}${project}${instance}"
storageAccountName="st${owner}${project}${instance}"
blobContainerName="${project}tfstates"
servicePrincipalName="sp${owner}${project}${instance}"

# Resource group
echo "Checking shared resource group [${resourceGroupName}]..."
resourceGroup=$(az group show \
  --name $resourceGroupName \
  2> /dev/null)

if [[ $resourceGroup == "" ]]; then
  echo " -> Shared resource group does not exist. Creating..."

  resourceGroup=$(az group create \
    --name $resourceGroupName \
    --location $location)

  echo -e " -> Shared resource group is created successfully.\n"
else
  echo -e " -> Shared resource group already exists.\n"
fi

# Key vault
echo "Checking shared key vault [${keyVaultName}]..."
keyVault=$(az keyvault show \
  --resource-group $resourceGroupName \
  --name $keyVaultName \
  2> /dev/null)
if [[ $keyVault == "" ]]; then
  echo " -> Shared key vault does not exist. Creating..."

  keyVault=$(az keyvault create \
    --resource-group $resourceGroupName \
    --name $keyVaultName \
    --location $location)

  echo -e " -> Shared key vault is created successfully.\n"
else
  echo -e " -> Shared key vault already exists.\n"
fi

# Storage account
echo "Checking shared storage account [${storageAccountName}]..."
storageAccount=$(az storage account show \
    --resource-group $resourceGroupName \
    --name $storageAccountName \
  2> /dev/null)

if [[ $storageAccount == "" ]]; then
  echo " -> Shared storage account does not exist. Creating..."

  storageAccount=$(az storage account create \
    --resource-group $resourceGroupName \
    --name $storageAccountName \
    --sku "Standard_LRS" \
    --allow-blob-public-access true \
    --encryption-services "blob")

  echo -e " -> Shared storage account is created successfully.\n"
else
  echo -e " -> Shared storage account already exists.\n"
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

# Service principal
echo "Checking service principal [${servicePrincipalName}]..."
subscriptionId=$(az account show | jq -r .id)

servicePrincipal=$(az ad sp show \
  --id $servicePrincipalName \
  2> /dev/null)
if [[ $servicePrincipal == "" ]]; then
  echo " -> Service principal does not exist. Creating..."

  servicePrincipal=$(az ad sp create-for-rbac \
    --name $servicePrincipalName \
    --role owner \
    --scopes "/subscriptions/${subscriptionId}" \
    --output json)

  echo -e " -> Service principal is created successfully.\n"

  servicePrincipalAppId=$(echo $servicePrincipal | jq -r .appId)
  servicePrincipalSecret=$(echo $servicePrincipal | jq -r .password)

  # Store service principal credentials in key vault
  echo -e " -> Storing service principal credentials into key vault."
  az keyvault secret set \
    --vault-name $keyVaultName \
    --name "${servicePrincipalName}-appid" \
    --value $servicePrincipalAppId \
    > /dev/null

  az keyvault secret set \
    --vault-name $keyVaultName \
    --name "${servicePrincipalName}-secret" \
    --value $servicePrincipalSecret \
    > /dev/null

  echo -e " -> Service principal credentials are stored successfully.\n"
else
  echo -e " -> Service principal already exists.\n"
fi


