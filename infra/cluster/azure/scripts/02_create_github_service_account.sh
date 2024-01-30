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
mainResourceGroupName="rg${project}main${instance}"
mainAksResourceName="aks${project}main${instance}"
mainKeyVaultName="kv${project}main${instance}"
servicePrincipalName="sp${project}${instance}"

# Resource group
echo "Checking main resource group [${mainResourceGroupName}]..."
resourceGroup=$(az group show \
  --name $mainResourceGroupName \
  2> /dev/null)

if [[ $resourceGroup == "" ]]; then
  echo -e " -> Main resource group does not exist! Create it first.\n"
  exit 1
else
  echo -e " -> Main resource groupd already exists.\n"
fi

# Key vault
echo "Checking main key vault [${mainKeyVaultName}]..."
keyVault=$(az keyvault show \
  --resource-group $mainResourceGroupName \
  --name $mainKeyVaultName \
  2> /dev/null)
if [[ $keyVault == "" ]]; then
  echo " -> Main key vault does not exist. Creating..."

  keyVault=$(az keyvault create \
    --resource-group $mainResourceGroupName \
    --name $mainKeyVaultName \
    --location $location)

  echo -e " -> Main key vault is created successfully.\n"
else
  echo -e " -> Main key vault already exists.\n"
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
    --role "Contributor" \
    --scopes "/subscriptions/${subscriptionId}/resourceGroups/${mainResourceGroupName}/providers/Microsoft.ContainerService/managedClusters/${mainAksResourceName}" \
    --output json)

  echo -e " -> Service principal is created successfully.\n"

  servicePrincipalAppId=$(echo $servicePrincipal | jq -r .appId)
  servicePrincipalSecret=$(echo $servicePrincipal | jq -r .password)

  # Store service principal credentials in key vault
  echo -e " -> Storing service principal credentials into key vault."
  az keyvault secret set \
    --vault-name $mainKeyVaultName \
    --name "${servicePrincipalName}-appid" \
    --value $servicePrincipalAppId \
    > /dev/null

  az keyvault secret set \
    --vault-name $mainKeyVaultName \
    --name "${servicePrincipalName}-secret" \
    --value $servicePrincipalSecret \
    > /dev/null

  echo -e " -> Service principal credentials are stored successfully.\n"
else
  echo -e " -> Service principal already exists.\n"
fi


