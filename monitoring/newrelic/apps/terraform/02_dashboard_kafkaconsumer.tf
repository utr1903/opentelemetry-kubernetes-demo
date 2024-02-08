#################
### Dashboard ###
#################

# Dashboard
resource "newrelic_one_dashboard" "kafkaconsumer" {
  name = "OTel Playground - Golang - Kafka Consumer"

  ###########################
  ### Runtime Performance ###
  ###########################

  # Golang
  dynamic "page" {
    for_each = var.language == "golang" ? [1] : []

    content {
      name = "Runtime Performance"

      # Go Routines
      widget_markdown {
        title  = "Go routines"
        column = 1
        row    = 1
        width  = 4
        height = 3

        text = "## Go routines\n\nThe following metric is considered:\n\n- Number of goroutines that currently exist\n -   `process.runtime.go.goroutines`"
      }

      # Average number of Go routines across all instances
      widget_billboard {
        title  = "Average number of Go routines across all instances"
        column = 5
        row    = 1
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.goroutines`) AS `Routines` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer'"
        }
      }

      # Average number of Go routines per instance
      widget_bar {
        title  = "Average number of Go routines per instance"
        column = 9
        row    = 1
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.goroutines`) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id"
        }
      }

      # Average number of Go routines across all instances
      widget_line {
        title  = "Average number of Go routines across all instances"
        column = 1
        row    = 4
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.goroutines`) AS `Routines` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' TIMESERIES"
        }
      }

      # Average number of Go routines per instance
      widget_line {
        title  = "Average number of Go routines per instance"
        column = 7
        row    = 4
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.goroutines`) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id TIMESERIES"
        }
      }

      # Garbage collection cycles
      widget_markdown {
        title  = "Garbage collection cycles"
        column = 1
        row    = 7
        width  = 4
        height = 3

        text = "## Garbage collection cycles\n\nThe following metric is considered:\n\n- Number of completed garbage collection cycles\n   - `process.runtime.go.gc.count`"
      }

      # Average number of garbage collection cycle across all instances
      widget_billboard {
        title  = "Average number of garbage collection cycle across all instances"
        column = 5
        row    = 7
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.gc.count`) AS `Routines` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer'"
        }
      }

      # Average number of garbage collection cycle per instance
      widget_bar {
        title  = "Average number of garbage collection cycle per instance"
        column = 9
        row    = 7
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.gc.count`) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id"
        }
      }

      # Average number of garbage collection cycle across all instances
      widget_line {
        title  = "Average number of garbage collection cycle across all instances"
        column = 1
        row    = 10
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.gc.count`) AS `Routines` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' TIMESERIES"
        }
      }

      # Average number of Go routines per instance
      widget_line {
        title  = "Average number of Go routines per instance"
        column = 7
        row    = 10
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.gc.count`) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id TIMESERIES"
        }
      }

      # Memory objects
      widget_markdown {
        title  = "Memory objects"
        column = 1
        row    = 13
        width  = 4
        height = 3

        text = "## Memory objects\n\nThe following metrics are considered:\n\n- Number of allocated heap objects\n   - `process.runtime.go.mem.heap_objects`\n- Number of live objects is the number of cumulative Mallocs - Frees\n   - `process.runtime.go.mem.live_objects`"
      }

      # Average number of memory objects across all instances
      widget_billboard {
        title  = "Average number of memory objects across all instances"
        column = 5
        row    = 13
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_objects`) AS `Heap`, average(`process.runtime.go.mem.live_objects`) AS `Live` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer'"
        }
      }

      # Average number of memory objects per instance
      widget_bar {
        title  = "Average number of memory objects per instance"
        column = 9
        row    = 13
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_objects`) AS `Heap`, average(`process.runtime.go.mem.live_objects`) AS `Live` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id"
        }
      }

      # Average number of memory objects across all instances
      widget_area {
        title  = "Average number of memory objects across all instances"
        column = 1
        row    = 16
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_objects`) AS `Heap`, average(`process.runtime.go.mem.live_objects`) AS `Live` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' TIMESERIES"
        }
      }

      # Average number of memory objects per instance
      widget_line {
        title  = "Average number of memory objects per instance"
        column = 7
        row    = 16
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_objects`) AS `Heap`, average(`process.runtime.go.mem.live_objects`) AS `Live` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id TIMESERIES"
        }
      }

      # Memory consumption (bytes)
      widget_markdown {
        title  = "Memory consumption (bytes)"
        column = 1
        row    = 19
        width  = 4
        height = 3

        text = "## Memory consumption\n\nThe following metrics are considered:\n\n- Bytes of allocated heap objects\n   - `process.runtime.go.mem.heap_alloc`\n- Bytes in idle (unused) spans\n   - `process.runtime.go.mem.heap_idle`\n- Bytes in in-use spans\n   - `process.runtime.go.mem.heap_inuse`\n- Bytes of idle spans whose physical memory has been returned to the OS\n   - `process.runtime.go.mem.heap_released`\n- Bytes of heap memory obtained from the OS\n   - `process.runtime.go.mem.heap_sys`"
      }

      # Average memory consumption across all instances (bytes)
      widget_billboard {
        title  = "Average memory consumption across all instances (bytes)"
        column = 5
        row    = 19
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_alloc`) AS `heap_alloc`, average(`process.runtime.go.mem.heap_idle`) AS `heap_idle`, average(`process.runtime.go.mem.heap_inuse`) AS `heap_inuse`, average(`process.runtime.go.mem.heap_released`) AS `heap_released`, average(`process.runtime.go.mem.heap_sys`) AS `heap_sys` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer'"
        }
      }

      # Average memory consumption per instance (bytes)
      widget_bar {
        title  = "Average memory consumption per instance (bytes)"
        column = 9
        row    = 19
        width  = 4
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_alloc`) AS `heap_alloc`, average(`process.runtime.go.mem.heap_idle`) AS `heap_idle`, average(`process.runtime.go.mem.heap_inuse`) AS `heap_inuse`, average(`process.runtime.go.mem.heap_released`) AS `heap_released`, average(`process.runtime.go.mem.heap_sys`) AS `heap_sys` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id"
        }
      }

      # Average memory consumption across all instances (bytes)
      widget_area {
        title  = "Average memory consumption across all instances (bytes)"
        column = 1
        row    = 22
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_alloc`) AS `heap_alloc`, average(`process.runtime.go.mem.heap_idle`) AS `heap_idle`, average(`process.runtime.go.mem.heap_inuse`) AS `heap_inuse`, average(`process.runtime.go.mem.heap_released`) AS `heap_released`, average(`process.runtime.go.mem.heap_sys`) AS `heap_sys` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' TIMESERIES"
        }
      }

      # Average memory consumption per instance (bytes)
      widget_line {
        title  = "Average memory consumption per instance (bytes)"
        column = 7
        row    = 22
        width  = 6
        height = 3

        nrql_query {
          account_id = var.NEW_RELIC_ACCOUNT_ID
          query      = "FROM Metric SELECT average(`process.runtime.go.mem.heap_alloc`) AS `heap_alloc`, average(`process.runtime.go.mem.heap_idle`) AS `heap_idle`, average(`process.runtime.go.mem.heap_inuse`) AS `heap_inuse`, average(`process.runtime.go.mem.heap_released`) AS `heap_released`, average(`process.runtime.go.mem.heap_sys`) AS `heap_sys` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET service.instance.id TIMESERIES"
        }
      }
    }
  }

  #########################################
  ### Application Performance (Metrics) ###
  #########################################
  page {
    name = "Application Performance (Metrics)"

    # Golden Signals
    widget_markdown {
      title  = ""
      column = 1
      row    = 1
      width  = 3
      height = 3

      text = "## Application Performance\n\nThis page is dedicated for the application golden signals retrieved from the metrics.\n\n- Latency\n- Throughput\n- Error Rate"
    }

    # Average latency across all instances (ms)
    widget_billboard {
      title  = "Average latency across all instances (ms)"
      column = 4
      row    = 1
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(messaging.receive.duration) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer'"
      }
    }

    # Total throughput across all instances (rpm)
    widget_billboard {
      title  = "Total throughput across all instances (rpm)"
      column = 7
      row    = 1
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(messaging.receive.duration), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer'"
      }
    }

    # Average error rate across all instances (%)
    widget_billboard {
      title  = "Average error rate across all instances (%)"
      column = 10
      row    = 1
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(messaging.receive.duration), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND `error.type` IS NOT NULL)/count(messaging.receive.duration)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer'"
      }
    }

    # Latency
    widget_markdown {
      title  = ""
      column = 1
      row    = 4
      width  = 3
      height = 3

      text = "## Latency\n\nLatency is monitored per the metric `messaging.receive.duration` which represents a histogram.\n\nIt corresponds to the aggregated consume time of the Kafka consumer.\n\nMoreover, the detailed performance can be investigated according to the topics, error types, instances etc."
    }

    # Average latency per Kafka topic across all instances (ms)
    widget_billboard {
      title  = "Average latency per Kafka topic across all instances (ms)"
      column = 4
      row    = 4
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(`messaging.receive.duration`) AS `Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET `messaging.destination.name`"
      }
    }

    # Average latency per error type across all instances (ms)
    widget_bar {
      title  = "Average latency per error type across all instances (ms)"
      column = 7
      row    = 4
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(`messaging.receive.duration`) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND `error.type`"
      }
    }

    # Average latency across all instances (ms)
    widget_line {
      title  = "Average latency across all instances (ms)"
      column = 1
      row    = 7
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(`messaging.receive.duration`) AS `Overall Latency` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' TIMESERIES"
      }
    }

    # Average latency per instance (ms)
    widget_line {
      title  = "Average latency per instance (ms)"
      column = 7
      row    = 7
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT average(`messaging.receive.duration`) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Throughput
    widget_markdown {
      title  = ""
      column = 1
      row    = 10
      width  = 3
      height = 3

      text = "## Throughput\n\nThroughput is monitored per the rate of change in the metric `messaging.receive.duration` in format of request per minute.\n\nIt corresponds to the aggregated amount of messages which are processed by the Kafka consumer in a minute.\n\nMoreover, the detailed performance can be investigated according to the topics, error types, instances etc."
    }

    # Total throughput per Kafka topic across all instances (rpm)
    widget_billboard {
      title  = "Total throughput per Kafka topic across all instances (rpm)"
      column = 4
      row    = 10
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(`messaging.receive.duration`), 1 minute) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET `messaging.destination.name`"
      }
    }

    # Total throughput per error type across all instances (rpm)
    widget_bar {
      title  = "Total throughput per error type across all instances (rpm)"
      column = 7
      row    = 10
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(`messaging.receive.duration`), 1 minute) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET `error.type`"
      }
    }

    # Total throughput across all instances (rpm)
    widget_line {
      title  = "Total throughput across all instances (rpm)"
      column = 1
      row    = 13
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(`messaging.receive.duration`), 1 minute) AS `Overall Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' TIMESERIES"
      }
    }

    # Average throughput per instance (rpm)
    widget_line {
      title  = "Average throughput per instance (rpm)"
      column = 7
      row    = 13
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT rate(count(`messaging.receive.duration`), 1 minute) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET k8s.pod.name TIMESERIES"
      }
    }

    # Error rate
    widget_markdown {
      title  = ""
      column = 1
      row    = 16
      width  = 3
      height = 3

      text = "## Error rate\n\nError rate is monitored per the metric `messaging.receive.duration` which ended with an error.\n\nIt corresponds to the ratio of the aggregated amount of consumed messages which have an error in compared to all consumed messages.\n\nMoreover, the detailed performance can be investigated according to the the topics, error types, instances etc."
    }

    # Average error rate per Kafka topic across all instances (%)
    widget_billboard {
      title  = "Average error rate per Kafka topic across all instances (%)"
      column = 4
      row    = 16
      width  = 3
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(`messaging.receive.duration`), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND numeric(`http.response.status_code`) >= 500)/count(`messaging.receive.duration`)*100 WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET `messaging.destination.name`"
      }
    }

    # Average error rate per error types across all instances (%)
    widget_bar {
      title  = "Error rate per error types across all instances (%)"
      column = 7
      row    = 16
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(`messaging.receive.duration`), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND `error.types` IS NOT NULL)/count(`messaging.receive.duration`)*100 WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET `error.type`"
      }
    }

    # Average error rate across all instances (%)
    widget_line {
      title  = "Average error rate across all instances (%)"
      column = 1
      row    = 19
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(`messaging.receive.duration`), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND `error.type` IS NOT NULL)/count(`messaging.receive.duration`)*100 AS `Overall Error Rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' TIMESERIES"
      }
    }

    # Average error rate per instance (%)
    widget_line {
      title  = "Error rate per instance (%)"
      column = 7
      row    = 19
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Metric SELECT filter(count(`messaging.receive.duration`), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND `error.type` IS NOT NULL)/count(`messaging.receive.duration`)*100 WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' FACET k8s.pod.name TIMESERIES"
      }
    }
  }

  #######################################
  ### Application Performance (Spans) ###
  #######################################
  page {
    name = "Application Performance (Spans)"

    # Overall server performance
    widget_markdown {
      title  = "Overall server performance"
      column = 1
      row    = 1
      width  = 3
      height = 3

      text = "## HTTP Server"
    }

    # Average web response time (ms)
    widget_line {
      title  = "Average web response time (ms)"
      column = 4
      row    = 1
      width  = 9
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT average(duration.ms) AS `Response time` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'server' TIMESERIES"
      }
    }

    # Average web throughput (rpm)
    widget_line {
      title  = "Average web throughput (rpm)"
      column = 1
      row    = 4
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT rate(count(*), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'server' TIMESERIES"
      }
    }

    # Average error rate (%)
    widget_line {
      title  = "Average error rate (%)"
      column = 7
      row    = 4
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT filter(count(*), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND otel.status_code = 'ERROR')/count(*)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'server' TIMESERIES"
      }
    }

    # Database performace
    widget_markdown {
      title  = "Database performace"
      column = 1
      row    = 7
      width  = 3
      height = 3

      text = "## Database"
    }

    # Average database time (ms)
    widget_line {
      title  = "Average database time (ms)"
      column = 4
      row    = 7
      width  = 9
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT average(duration.ms) AS `DB time` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'client' AND server.address = 'mysql.golang.svc.cluster.local' TIMESERIES"
      }
    }

    # Database throughput (rpm)
    widget_line {
      title  = "Database throughput (rpm)"
      column = 1
      row    = 10
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT rate(count(*), 1 minute) AS `Throughput` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'client' AND server.address = 'mysql.golang.svc.cluster.local' TIMESERIES"
      }
    }

    # Database error rate (%)
    widget_line {
      title  = "Database error rate (%)"
      column = 7
      row    = 10
      width  = 6
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT filter(count(*), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND otel.status_code = 'ERROR')/count(*)*100 AS `Error rate` WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'client' AND server.address = 'mysql.golang.svc.cluster.local' TIMESERIES"
      }
    }

    # Max database operation latency (ms)
    widget_bar {
      title  = "Max database operation latency (ms)"
      column = 1
      row    = 13
      width  = 4
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT max(duration.ms) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'client' AND server.address = 'mysql.golang.svc.cluster.local' FACET db.name, db.table, db.operation"
      }
    }

    # Database operation throughput (rpm)
    widget_bar {
      title  = "Max database operation throughput (rpm)"
      column = 5
      row    = 13
      width  = 4
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT rate(count(*), 1 minute) WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'client' AND server.address = 'mysql.golang.svc.cluster.local' FACET db.name, db.table, db.operation"
      }
    }

    # Database operation error rate (%)
    widget_bar {
      title  = "Average database error rate (%)"
      column = 9
      row    = 13
      width  = 4
      height = 3

      nrql_query {
        account_id = var.NEW_RELIC_ACCOUNT_ID
        query      = "FROM Span SELECT filter(count(*), WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND otel.status_code = 'ERROR')/count(*)*100 WHERE instrumentation.provider = 'opentelemetry' AND k8s.cluster.name = '${var.cluster_name}' AND service.name = 'kafkaconsumer' AND span.kind = 'client' AND server.address = 'mysql.golang.svc.cluster.local' FACET db.name, db.table, db.operation"
      }
    }
  }
}
