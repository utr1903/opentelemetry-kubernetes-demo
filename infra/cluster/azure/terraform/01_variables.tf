#################
### Variables ###
#################

# Resource group name of AKS
variable "aks_resource_group_name" {
  type = string
}

# Resource name of AKS
variable "aks_resource_name" {
  type = string
}

# Resource group name of AKS nodepool
variable "aks_nodepool_resource_name" {
  type = string
}

# Kubernetes version of AKS
variable "aks_version" {
  type = string
}

# Datacenter location of AKS
variable "location" {
  type    = string
  default = "westeurope"
}
