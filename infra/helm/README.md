# Application deployment

In order to run the environment properly, the applications should be deployed in the following order.

1. [cert-manager](/infra/helm/cert-manager)
2. [oteloperator](/infra/helm/oteloperator/)
3. [otelcollector](/infra/helm/otelcollector/)
4. [kafka](/infra/helm/kafka/) & [mysql](/infra/helm/mysql/)
5. [httpserver](/infra/helm/httpserver/) & [kafkaconsumer](/infra/helm/kafkaconsumer/)
6. [simulator](/infra/helm/simulator/)

## Dependencies

- `oteloperator` requires `cert-manager` for certificates.
- `otelcollector` requires `oteloperator` for CRDs and Target Allocator.
- `httpserver` requires `mysql` & `otelcollector`.
- `kafkaconsumer` requires `mysql`, `kafka` & `otelcollector`.
- `simulator` requires `kafka`, `kafkaconsumer`, `httpserver` & `otelcollector`.

## Docker images

Our own applications (`httpserver`, `kafkaconsumer` & `simulator`) need to be built and pushed to a container registry for which the Github container registry (`ghcr.io`) is chosen. This step is automated by the Github workflow [`docker_build_push.yaml`](/.github/workflows/docker_build_push.yaml)

The workflow is designed to be run on demand and therefore it is to be run manually. It has 2 input parameters: `language` & `application`.

```
application
- httpserver
- kafkaconsumer
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

### Build

The workflow will refer to the main Dockerfile (`/apps/<language>/Dockerfile`) and inject the necessary application name into it. Depending on that, the necessary application image will be built for `linux/amd64` and `linux/arm64` architectures.

### Push

After the build step is successfully completed, the image will be pushed to the Github container registry. Every image will have the following naming convention:

`ghcr.io/${{ github.actor }}/${{ secrets.PROJECT }}-${{ inputs.application }}-${{ inputs.language }}:${{ github.sha }}`.

Example:
`ghcr.io/utr1903/myproj-httpserver-golang:xxx`

Moreover, the last pushed image will always be tagged as `latest` so that the Kubernetes cluster can pull the newer version of the image when a new Helm chart is deployed.

**IMPORTANT:** When you push an image to the Github container registry, the image repository will be PRIVATE by default. You have to manually change the visibility to PUBLIC because otherwise the Kubernetes cluster will not be able to pull the image!

## Helm deployment

Every application should be deployed in the order mentioned at the beginning of this document. In order to deploy any application, the Github workflow [`helm_deploy.yaml`](/.github/workflows/helm_deploy.yaml) is to be used.

Like the `docker_build_push.yaml`, this is also to be run on demand and therefore manually. It has 2 input parameters: `chart` & `language`.

```
chart
- cert-manager
- httpserver
- kafka
- kafkaconsumer
- mysql
- otelcollector
- oteloperator
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

The deployment of the charts `cert-manager`, `kafka`, `mysql`, `oteloperator` & `otelcollector` should leave the `language` input empty! The are using remote Helm charts for deployment.

On the other hand, `language` is required for the charts `httpserver`, `kafkaconsumer` & `simulator` because they are using local Helm charts!
