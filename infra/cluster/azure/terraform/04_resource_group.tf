### Resource Group ###

# Resource Group
resource "azurerm_resource_group" "platform" {
  name     = var.aks_resource_group_name
  location = var.location
}
