# Application build

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

## Local development

It is quite an effort to build and push an image to the Github container registry. One does not simply want to wait that long to test small changes. In order to test your application faster, there is the local development possibility.

**Remark:** If you do not want to create a cluster in a cloud provider, you can always use a `kind` cluster for development purposes. Check out how to deploy one [here](/infra/cluster/kind).

### Build

To build the images for your version of the application, run the script [build_push_local.sh](/apps/build_push_local.sh) as follows:

```shell
bash build_push_local.sh \
  --docker-username <DOCKERHUB_USERNAME> \
  --platform <YOUR_MACHINES_ARCHITECTURE> \ # amd64 or arm64
  --language <PROGRAMMING_LANGUAGE_OF_APPLICATION> \ # golang, java etc.
  --project <YOUR_PREFERRED_PROJECT_NAME>
```

Example:

```shell
bash build_push_local.sh --docker-username utr1903 --platform arm64 --language golang --project myproj
```

The image will have the following naming convention:
`docker.io/${project}-${application}-${language}:latest`

Example:
`docker.io/utr1903/myproj-httpserver-golang:xxx`

### Push

The built image will be pushed to your Dockerhub container registry right away.

**IMPORTANT:** You have to be logged in to your Dockerhub account!
