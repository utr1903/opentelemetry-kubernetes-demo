# Azure

## Setting up baseline

For automating Azure Kubernetes Service (AKS) deployment and start/stop, we will be using Github workflows and Terraform. In order for Github workflow to deploy anything on our Azure account, it requires access rights for which we will be using an Azure service principal. Moreover, the Terraform deployment requires to store the state of the deployment and for that we will be needing a blob container in a storage account.

This pre-setup is called the baseline and can be prepared with the script [`00_create_baseline_resources.sh`](/infra/cluster/azure/scripts/00_create_baseline_resources.sh). It is important to note that this script should be run by a user (or a sevice principal) who has `Owner` rights on the subscription level!

The following Azure resources will be deployed:

1. Resource group (to group all baseline resources)
2. Key vault (to store service principal credentials)
3. Storage account & blob container (to store Terraform state)
4. Service principal (to run Github workflows)
