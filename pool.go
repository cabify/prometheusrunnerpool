package prometheusrunnerpool

import (
	"github.com/cabify/runnerpool"
	"github.com/prometheus/client_golang/prometheus"
)

// Observe binds a pool to a prometheus gauge in the provided registerer,
// and starts observing it on every prometheus scraping interval
// If registering fails, it panics
func Observe(registry prometheus.Registerer, pool runnerpool.Pool, poolName string) {
	collectors := make([]prometheus.Collector, len(metrics))
	for i, m := range metrics {
		collectors[i] = newGauge(pool, poolName, m)
	}
	registry.MustRegister(collectors...)
}

type metric struct {
	name string
	help string
	stat func(s runnerpool.Stats) int32
}

// metrics maps the metric name to the value of the runnerpool.Stats (and vice-versa)
var metrics = []metric{
	{
		name: "max_workers_total",
		help: "Maximum amount of workers configured in the pool",
		stat: func(s runnerpool.Stats) int32 { return s.MaxWorkers },
	},
	{
		name: "workers_total",
		help: "Total number of workers in the pool",
		stat: func(s runnerpool.Stats) int32 { return s.Workers },
	},
	{
		name: "acquired_total",
		help: "Total number of workers acquired in the pool, they may or may not be running something",
		stat: func(s runnerpool.Stats) int32 { return s.Acquired },
	},
	{
		name: "running_total",
		help: "Total number of workers running a function in the pool",
		stat: func(s runnerpool.Stats) int32 { return s.Running },
	},
}

func newGauge(pool runnerpool.Pool, poolName string, m metric) prometheus.GaugeFunc {
	opts := prometheus.GaugeOpts{
		Namespace: "runnerpool",
		Subsystem: "stats",
		Name:      m.name,
		Help:      m.help,
		ConstLabels: map[string]string{
			"pool_name": poolName,
		},
	}
	return prometheus.NewGaugeFunc(opts, func() float64 {
		return float64(m.stat(pool.Stats()))
	})
}
