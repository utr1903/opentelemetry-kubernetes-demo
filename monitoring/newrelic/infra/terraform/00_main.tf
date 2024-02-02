#################
### New Relic ###
#################

# Kubernetes infra monitoring
module "nrotelk8s" {
  source = "github.com/newrelic-experimental/monitoring-kubernetes-with-opentelemetry.git?ref=newrelic-monitoring-0.2.0/monitoring/terraform"

  NEW_RELIC_ACCOUNT_ID = var.NEW_RELIC_ACCOUNT_ID
  NEW_RELIC_API_KEY    = var.NEW_RELIC_API_KEY
  NEW_RELIC_REGION     = var.NEW_RELIC_REGION
  cluster_name         = var.cluster_name
}
