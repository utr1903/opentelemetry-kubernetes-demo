# Kind

For local development purposes, `kind` is a magnificent tool! Summarized, it is providing you a Kubernetes cluster running in Docker containers (as in each container represents a node). Thereby, you can mock an actual cluster locally on your machine through the container runtime.

You can deploy a 4-node `kind` cluster as follows:

```shell
bash 00_deploy_cluster.sh --project myproj --instance 001 --k8s-version 1.28.0
```

and clean it up as follows:

```shell
bash 00_deploy_cluster.sh --project myproj --instance 001 --k8s-version 1.28.0 --destroy
```

This will not only deploy a cluster on your machine but also automatically switches the Kubernetes context to itself (meaning that your `kubectl` commands will be executed against your `kind` cluster)!
