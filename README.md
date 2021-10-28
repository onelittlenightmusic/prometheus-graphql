# prometheus-graphql

Simple GraphQL generator which generates GraphQL endpoint from Prometheus server.

<img src="./prometheus-graphql.png" width=400px>

## TL;DR

Start GraphQL server.

```sh
go run app/server.go
```

Build GraphQL server.

```sh
go build -o bin/server app/server.go
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

# Configure

Create a `config/config.yaml` depending on your usage.

```config.yaml
spec:
  prometheusAddress: http://192.168.0.71:9090/
```

Mount this file on `./config/` directory or `/config/` directory in container.

# Another run type

## Docker

```sh
docker run -d --name graphql -p 8080:8080 --volume $(PWD)/config:/config hiroyukiosaki/prometheus-graphql:latest 
```

or configure `docker-compose.yaml` and run this command.

```sh
docker-compose up -d
```

And access `http://localhost:8080`.

## Kubernetes

Just run 

```sh
kubectl apply -f k8s
```

```sh
kubectl port-forward svc/prom-graphql-svc 8080:8080
```

And access `http://localhost:8080`.