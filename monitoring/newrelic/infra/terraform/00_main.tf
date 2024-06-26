#################
### New Relic ###
#################

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">=3.89.0"
    }
  }

  backend "azurerm" {}
}

# Configure the Azure Provider
provider "azurerm" {
  features {}
}

# Kubernetes infra monitoring
module "nrotelk8s" {
  source = "github.com/newrelic-experimental/monitoring-kubernetes-with-opentelemetry.git?ref=newrelic-monitoring-0.3.1/monitoring/terraform"

  NEW_RELIC_ACCOUNT_ID = var.NEW_RELIC_ACCOUNT_ID
  NEW_RELIC_API_KEY    = var.NEW_RELIC_API_KEY
  NEW_RELIC_REGION     = var.NEW_RELIC_REGION
  cluster_name         = var.cluster_name
}
