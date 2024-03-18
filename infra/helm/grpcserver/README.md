# gRPC server

The `grpcserver` is running within the `<language>` namespace. It uses

- Redis cache

within the `ops` namespace.

The corresponding `GET` and `DELETE` producure works as follows:

1. `simulator` calls `grpcserver`
2. `httpserver` gets/deletes from/in `redis`

![workflow](/media/grpcserver_workflow.png)
