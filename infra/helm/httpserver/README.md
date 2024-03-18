# HTTP server

The `httpserver` is running within the `<language>` namespace. It uses

- Redis cache
- MySQL database

within the `ops` namespace.

The corresponding `GET` and `DELETE` producure works as follows:

1. `simulator` calls `httpserver`
2. `httpserver` gets the boolean `increase.latency` flag from `redis`
   - _periodic intentional error controlled by `latencymanager`_
3. `httpserver` gets/deletes from/in `mysql`

![workflow](/media/httpserver_workflow.png)
