### Kubernetes Cluster ###

# Kubernetes Cluster
resource "azurerm_kubernetes_cluster" "platform" {
  resource_group_name = azurerm_resource_group.platform.name
  location            = azurerm_resource_group.platform.location
  name                = var.aks_resource_name

  dns_prefix         = "${var.aks_resource_name}-${azurerm_resource_group.platform.name}"
  kubernetes_version = var.aks_version

  node_resource_group = var.aks_nodepool_resource_name

  default_node_pool {
    name    = "system"
    vm_size = "Standard_D2_v2"

    node_labels = {
      nodePoolName = "system"
    }

    enable_auto_scaling = true
    node_count          = 1
    min_count           = 1
    max_count           = 1
  }

  identity {
    type = "SystemAssigned"
  }
}

# Kubernetes Nodepool - General usage
resource "azurerm_kubernetes_cluster_node_pool" "general" {
  name                  = "general"
  kubernetes_cluster_id = azurerm_kubernetes_cluster.platform.id
  vm_size               = "Standard_D2_v2"

  orchestrator_version = var.aks_version

  node_labels = {
    nodePoolName = "general"
  }

  enable_auto_scaling = true
  node_count          = 3
  min_count           = 3
  max_count           = 3
}
