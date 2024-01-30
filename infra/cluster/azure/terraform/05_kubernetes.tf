### Kubernetes Cluster ###

# Kubernetes Cluster
resource "azurerm_kubernetes_cluster" "platform" {
  resource_group_name = azurerm_resource_group.platform.name
  location            = azurerm_resource_group.platform.location
  name                = var.aks_resource_name

  dns_prefix         = "${var.aks_resource_name}-${azurerm_resource_group.platform.name}"
  kubernetes_version = var.aks_version

  node_resource_group = var.aks_nodepool_resource_name

  network_profile {
    network_plugin = "kubenet"
    network_policy = "calico"
    load_balancer_sku = "basic"
  }

  default_node_pool {
    name    = "system"
    vm_size = "Standard_D2_v2"

    node_labels = {
      nodePoolName = "system"
    }

    enable_auto_scaling = false
    node_count          = 4
  }

  identity {
    type = "SystemAssigned"
  }
}
