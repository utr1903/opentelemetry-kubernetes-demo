# Latency manger

The `latencymanager` is running within the `<language>` namespace. It uses

- Redis cache

within the `ops` namespace.

The corresponding `increase` and `decrease` producure works as follows:

1. `cronjob` triggers `latencymanager`
   - _either `increase` or `decrease`_
2. `latencymanager` sets the boolean `increase.latency` flag in `redis`
   - _periodic intentional error to manipulate `httpserver` latency_

![workflow](/media/latencymanager_workflow.png)

Moreover, it sends a deployment marker event to New Relic whenever the latency change has been made.
