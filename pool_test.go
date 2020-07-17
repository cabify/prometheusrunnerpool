package prometheusrunnerpool_test

import (
	"strings"
	"testing"

	"github.com/cabify/prometheusrunnerpool"

	"github.com/cabify/runnerpool"
	"github.com/cabify/runnerpool/runnerpoolmock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

const nl = "\n"

func TestObserver_Observe(t *testing.T) {
	t.Run("reports metrics for one pool", func(t *testing.T) {
		registry := prometheus.NewRegistry()

		poolMock := &runnerpoolmock.Pool{}
		poolMock.On("Stats").Return(runnerpool.Stats{
			MaxWorkers: 4,
			Workers:    3,
			Acquired:   2,
			Running:    1,
		})

		prometheusrunnerpool.Observe(registry, poolMock, "test_pool")

		reader := strings.NewReader(`` +
			`# HELP runnerpool_stats_max_workers_total Maximum amount of workers configured in the pool` + nl +
			`# TYPE runnerpool_stats_max_workers_total gauge` + nl +
			`runnerpool_stats_max_workers_total{pool_name="test_pool"} 4` + nl +

			`# HELP runnerpool_stats_workers_total Total number of workers in the pool` + nl +
			`# TYPE runnerpool_stats_workers_total gauge` + nl +
			`runnerpool_stats_workers_total{pool_name="test_pool"} 3` + nl +

			`# HELP runnerpool_stats_acquired_total Total number of workers acquired in the pool, they may or may not be running something` + nl +
			`# TYPE runnerpool_stats_acquired_total gauge` + nl +
			`runnerpool_stats_acquired_total{pool_name="test_pool"} 2` + nl +

			`# HELP runnerpool_stats_running_total Total number of workers running a function in the pool` + nl +
			`# TYPE runnerpool_stats_running_total gauge` + nl +
			`runnerpool_stats_running_total{pool_name="test_pool"} 1` + nl,
		)

		err := testutil.GatherAndCompare(registry, reader,
			"runnerpool_stats_max_workers_total",
			"runnerpool_stats_workers_total",
			"runnerpool_stats_acquired_total",
			"runnerpool_stats_running_total",
		)
		assert.NoError(t, err)
	})

	t.Run("reports metrics for two pools", func(t *testing.T) {
		registry := prometheus.NewRegistry()

		firstPoolMock := &runnerpoolmock.Pool{}
		firstPoolMock.On("Stats").Return(runnerpool.Stats{
			MaxWorkers: 4,
			Workers:    3,
			Acquired:   2,
			Running:    1,
		})
		secondPoolMock := &runnerpoolmock.Pool{}
		secondPoolMock.On("Stats").Return(runnerpool.Stats{
			MaxWorkers: 40,
			Workers:    30,
			Acquired:   20,
			Running:    10,
		})

		prometheusrunnerpool.Observe(registry, firstPoolMock, "first_pool")
		prometheusrunnerpool.Observe(registry, secondPoolMock, "second_pool")

		reader := strings.NewReader(`` +
			`# HELP runnerpool_stats_max_workers_total Maximum amount of workers configured in the pool` + nl +
			`# TYPE runnerpool_stats_max_workers_total gauge` + nl +
			`runnerpool_stats_max_workers_total{pool_name="first_pool"} 4` + nl +
			`runnerpool_stats_max_workers_total{pool_name="second_pool"} 40` + nl +

			`# HELP runnerpool_stats_workers_total Total number of workers in the pool` + nl +
			`# TYPE runnerpool_stats_workers_total gauge` + nl +
			`runnerpool_stats_workers_total{pool_name="first_pool"} 3` + nl +
			`runnerpool_stats_workers_total{pool_name="second_pool"} 30` + nl +

			`# HELP runnerpool_stats_acquired_total Total number of workers acquired in the pool, they may or may not be running something` + nl +
			`# TYPE runnerpool_stats_acquired_total gauge` + nl +
			`runnerpool_stats_acquired_total{pool_name="first_pool"} 2` + nl +
			`runnerpool_stats_acquired_total{pool_name="second_pool"} 20` + nl +

			`# HELP runnerpool_stats_running_total Total number of workers running a function in the pool` + nl +
			`# TYPE runnerpool_stats_running_total gauge` + nl +
			`runnerpool_stats_running_total{pool_name="first_pool"} 1` + nl +
			`runnerpool_stats_running_total{pool_name="second_pool"} 10` + nl,
		)

		err := testutil.GatherAndCompare(registry, reader,
			"runnerpool_stats_max_workers_total",
			"runnerpool_stats_workers_total",
			"runnerpool_stats_acquired_total",
			"runnerpool_stats_running_total",
		)
		assert.NoError(t, err)
	})
}
