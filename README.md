# prometheus-graphql

Simple GraphQL generator which generates GraphQL endpoint from Prometheus server.

## TL;DR

Start GraphQL server.

```sh
go run server.go
```

Access GraphQL Playground interface.

http://localhost:8080

## Requirement

- **Prometheus should run** on a reachable location.
  - Before running prometheus-graphql, please configure `config.yaml` to point your Prometheus location.

### Configuration: connect to Prometheus

```yaml
spec:
  prometheusAddress: http://your.prometheus.url:9090
```

## Query

```graphql
{
  label_values(label: "__name__")
}
```

```json
{
  "data": {
    "label_values": [
      "aggregator_openapi_v2_regeneration_count",
      "aggregator_openapi_v2_regeneration_duration",
      "aggregator_unavailable_apiservice",
      "aggregator_unavailable_apiservice_total",
      ...
    ]
  }
}
```

This query can be replaced with just ` { name_values }`.

---

```graphql
{
  labels
}
```

```json
{
  "data": {
    "labels": [
      "__name__",
      "access_mode",
      "apiservice",
      "build_date",
      "claim_namespace",
      "cluster_ip",
      "code",
      ...
    ]
  }
}
```

---

> Querying metrics

```graphql
{
  query_range(query: "avg_over_time(kube_node_status_allocatable[1m])") {
    metric
    values {
      value
      timestamp
    }
  }
}
```

```json
{
  "data": {
    "query_range": [
      {
        "metric": {
          "instance": "kube-state-metrics.kube-system.svc.cluster.local:8080",
          "job": "kube-state-metrics",
          "node": "master-node-1",
          "resource": "cpu",
          "unit": "core"
        },
        "values": [
          {
            "value": 4,
            "timestamp": 1621582576547
          },
          {
            "value": 4,
            "timestamp": 1621582696547
          },
          ...
```

This metrics can be fetched if you installed `kube-state-metrics` on your Kubernetes cluster.