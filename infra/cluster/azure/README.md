# Azure

## 01 - Setting up baseline

For provisioning Azure Kubernetes Service (AKS), we will be using Terraform which needs to store the state of the deployment somewhere. That storage is a blob container within a storage account which is supposed to be available in advance to the AKS provisioning.

This pre-setup is called the baseline and can be created with the script [`00_create_baseline_resources.sh`](/infra/cluster/azure/scripts/00_create_baseline_resources.sh). It is important to note that this script should be run by a user (or a sevice principal) who has `Owner` rights on the subscription level!

```shell
bash 00_create_baseline_resources.sh --project myproj --instance 001 --location westeurope
```

The following Azure resources will be deployed:

1. Resource group [`rg${project}base${instance}`] (to group all baseline resources)
2. Storage account [`st${project}base${instance}`] (for Terraform blob container)
3. Blob container [`${project}tfstates`] (to store Terraform state)

## 02 - Provisioning cluster

After the baseline components are created, the cluster and it's relevant resources can be provisioned. In order to do that, run the script [`01_deploy_cluster.sh`](/infra/cluster/azure/scripts/01_deploy_cluster.sh).

```shell
bash 01_deploy_cluster.sh --project myproj --instance 001 --location westeurope --k8s-version 1.28.0
```

**IMPORTANT**: The `project`, `instance` and `location` should be the same as the ones in the baseline!

The following Azure resources will be deployed:

1. Resource group [`rg${project}main${instance}`] (to group all main resources)
2. AKS [`aks${project}main${instance}`] (cluster itself)
3. AKS resource group [`rgaks${project}main${instance}`] (to group the cluster nodepools)
4. Key vault [`kv${project}main${instance}`] (to store service principal credentials)

## 03 - Creating service account for Github actions

In order for a Github workflow to start/stop or deploy Helm charts onto the AKS, it requires access rights for which we will be using an Azure service principal. That will be created and given the necessary rights per the script [`02_create_github_service_account.sh`](/infra/cluster/azure/scripts/02_create_github_service_account.sh).

Moreover, the same service principal will be used to run Terraform deployments which needs to store the state of the deployment in a backend. This backend will be the base storage account `st${project}base${instance}`. In order to create store that state in a blob container, the service principal will be given necessary rights on the storage account as well!

```shell
bash 02_create_github_service_account.sh --project myproj --instance 001 --location westeurope
```

**IMPORTANT**: The `project`, `instance` and `location` should be the same as the ones for the AKS!

The following Azure resources will be deployed:

1. Service principal [`sp${project}${instance}`] (to run Github workflows)
2. Key vault secrets (to store service principal credentials)

After this step is completed, we need to store the following parameters as Github secrets so that our Github workflows can talk to Azure successfully:

1. PROJECT (from the flag `--project`)
2. INSTANCE (from the flag `--instance`)
3. AZURE_TENANT_ID (your Azure tenant ID)
4. AZURE_SUBSCRIPTION_ID (your Azure subscription ID)
5. AZURE_SERVICE_PRINCIPAL_APP_ID (app ID of the service principal)
6. AZURE_SERVICE_PRINCIPAL_SECRET (secret of the service principal)

## 04 - Cleaning up

After we are done with the entire environment, we need to clean up everything we have created. In order to that, do the following sequentially:

First, we destroy the Terraform deployment:

```shell
bash 01_deploy_cluster.sh --project myproj --instance 001 --location westeurope --k8s-version 1.28.0 --destroy
```

Next, we delete the service principal (and the app registration behind it) and remove all of the baseline resources by running the clean up script [`03_cleanup_resources.sh`](/infra/cluster/azure/scripts/03_cleanup_resources.sh).

```shell
bash 03_cleanup_resources.sh --project myproj --instance 001 --location westeurope --destroy
```

**IMPORTANT**: The `project`, `instance` and `location` should be the same as the ones in the baseline and main!

Last, the key vault needs to be purged. Azure deletes the key vaults in a soft manner so that they can be recovered. In order to permanently remove a key vault, it has to be purged. That's something we need to do manually in the portal.
