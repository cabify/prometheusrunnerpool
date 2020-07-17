# prometheusrunnerpool

## Prometheus metrics for `github.com/cabify/runnerpool`

Usage:

In order to register a pool in the default prometheus registry:
```go
prometheusrunnerpool.Observe(prometheus.DefaultRegistry, pool, "http-requests")
```
