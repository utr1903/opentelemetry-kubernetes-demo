# Application deployment

In order to run the environment properly, the applications should be deployed in the following order.

1. [cert-manager](/infra/helm/cert-manager)
2. [oteloperator](/infra/helm/oteloperator/)
3. [otelcollector](/infra/helm/otelcollector/)
4. [kafka](/infra/helm/kafka/) & [mysql](/infra/helm/mysql/)
5. [httpserver](/infra/helm/httpserver/) & [kafkaconsumer](/infra/helm/kafkaconsumer/)
6. [simulator](/infra/helm/simulator/)
7. [latencymanager](/infra/helm/latencymanager/)

## Dependencies

- `oteloperator` requires `cert-manager` for certificates.
- `otelcollector` requires `oteloperator` for CRDs and Target Allocator.
- `httpserver` requires `redis`, `mysql` & `otelcollector`.
- `kafkaconsumer` requires `mysql`, `kafka` & `otelcollector`.
- `simulator` requires `kafka`, `kafkaconsumer`, `httpserver` & `otelcollector`.
- `latencymanager` requires `redis` & `otelcollector`.

## Helm deployment

Every application should be deployed in the order mentioned at the beginning of this document. In order to deploy any application, the Github workflow [`helm_deploy.yaml`](/.github/workflows/helm_deploy.yaml) is to be used.

Like the `docker_build_push.yaml`, this is also to be run on demand and therefore manually. It has 2 input parameters: `chart` & `language`.

```
chart
- cert-manager
- httpserver
- kafka
- kafkaconsumer
- latencymanager
- mysql
- otelcollector
- oteloperator
- redis
- simulator
```

```
language
- dotnet
- golang
- java
- javascript
- python
```

**IMPORTANT**

The deployment of the charts `cert-manager`, `kafka`, `mysql`, `oteloperator`, `otelcollector` & `redis` should leave the `language` input empty! They are using remote Helm charts for deployment.

On the other hand, `language` is required for the charts `httpserver`, `kafkaconsumer`,`latencymanager` & `simulator` because they are using local Helm charts!

## Local development

When you are developing locally, you wouldn't want to push your code to Github, wait for a runner to pick it up and deploy it to your cluster. In matter of fact, if you are working locally on a `kind` cluster, the Helm deployment of the Github workflow won't help you at all.

### Remote Helm charts

The deployment of the charts `cert-manager`, `oteloperator`, `kafka`, `mysql` & `redis` are fairly simple. Just change directory to their corresponding folders and run the `deploy.sh` script. It will automatically deploy the remote Helm charts onto your `kind` cluster with necessary configuration.

The `otelcollector` is the key component of this repository. This is a special implementation of various types of collectors and is the hub for all the telemetry data collected throughout the cluster. For detailed explanation of this Helm chart, refer to it's actual [repository](https://github.com/newrelic-experimental/monitoring-kubernetes-with-opentelemetry)!

You can deploy it as follows:

```shell
bash deploy.sh \
  --project <YOUR_PREFERRED_PROJECT_NAME> \
  --instance <YOUR_INSTANCE> \
  --cluster-type <YOUR_CLUSTER_TYPE> \ # aks, eks, gke, kind
  --newrelic-otlp-endpoint \ # https://otlp.nr-data.net or https://otlp.eu01.nr-data.net
  --newrelic-opsteam-license-key <YOUR_NEWRELIC_LICENSE_KEY>
```

Example:

```shell
bash deploy.sh --project myproj --instance 001 --cluster-type kind --newrelic-otlp-endpoint https://otlp.eu01.nr-data.net --newrelic-opsteam-license-key $NEWRELIC_LICENSE_KEY_OPSTEAM
```

### Local Helm charts

Our own applications `httpserver`, `kafkaconsumer`, `latencymanager` & `simulator` have their own local Helm templates and thereby are to be deployed with the local Helm configuration. You can run the `/infra/helm/${language}/deploy_local.sh` scripts in each application.

Example (`httpserver`):

```shell
bash deploy_local.sh \
  --docker-username <DOCKERHUB_USERNAME> \
  --language <PROGRAMMING_LANGUAGE_OF_APPLICATION> \ # golang, java etc.
  --project <YOUR_PREFERRED_PROJECT_NAME> \
  --instance <YOUR_INSTANCE>
```

Example:

```shell
bash deploy_local.sh --docker-username utr1903 --language golang --project myproj --instance 001
```
