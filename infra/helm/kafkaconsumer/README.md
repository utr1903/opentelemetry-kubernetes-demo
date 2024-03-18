# Kafka consumer

The `kafkaconsumer` is running within the `<language>` namespace. It uses

- Kafka broker
- Redis cache
- MySQL database

within the `ops` namespace.

The corresponding `CREATE` producure works as follows:
1. `simulator` publishes a message to `kafka`
2. `kafkaconsumer` consumes the message
3. `kafkaconsumer` stores a value in `redis`
4. `kafkaconsumer` stores a value in `mysql`

![workflow](/media/kafkaconsumer_workflow.png)
