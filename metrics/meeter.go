package collector

import (
	"context"
	"fmt"
	"github.com/ISE-SMILE/corral/api"
	"github.com/docker/docker/api/types/filters"
	"math/rand"
	"net/http"
	time "time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Metrics struct {
	nodes    []string
	ctx      context.Context
	interval time.Duration

	server      *http.Server
	nodeMetrics map[string]*NodeStats
	metrics     *api.Metrics
}

const dockerMetricsPort = 2376
const nodeMetricsPort = 9100

func nodeFields() []string {
	return []string{"node_memory_Active", "node_memory_MemAvailable", "node_cpu", "node_disk_io_now", "node_load1", "node_memory_MemFree", "node_memory_MemFree", "node_memory_MemTotal", "node_network_receive_bytes", "node_network_transmit_bytes"}
}

var header = map[string]string{
	"timestamp": "current timestamp",
	"HId":       "node name",

	"cpu_user":         "cpu used py user in percent",
	"cpu_system":       "cpu used py system in percent",
	"cpu_idle":         "cpu idle in percent",
	"cpu_iowait":       "cpu used py iowait in percent",
	"forks":            "number of forks",
	"disk":             "disc usage",
	"load":             "running load average (1 minute)",
	"mem_available":    "avaliable memory in bytes",
	"mem_total":        "total memory in bytes",
	"net_rev":          "bytes received",
	"net_tra":          "bytes transmitted",
	"cpu_user_delta":   "cpu user changed since last update",
	"cpu_system_delta": "cpu system changed since last update",
	"cpu_idle_delta":   "cpu idle changed since last update",
	"cpu_iowait_delta": "cpu iowait changed since last update",
	"forks_delta":      "forks changed since last update",
	"mem_usage":        "memory usage in bytes",

	"docker_running": "number of running docker containers",
	"docker_paused":  "number of paused docker containers",
	"docker_total":   "total number of docker containers",
}

func (m *Metrics) Setup(metrics *api.Metrics) error {
	if metrics != nil {

		metrics, err := api.CollectMetrics(header)
		if err != nil {
			return err
		}
		m.metrics = metrics

	} else {
		m.metrics = metrics
		for k, v := range header {
			err := m.metrics.AddField(k, v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *Metrics) Collect(duration time.Duration, filter filters.Args) error {
	m.interval = duration

	errors := make(chan error)

	nodeMetrics := make(chan NodeMetrics)
	dockerMetrics := make(chan DockerMetrics)
	for _, node := range m.nodes {
		go prometheusCollector(m.ctx, node, fmt.Sprintf("http://%s:%d/metrics", node, nodeMetricsPort),
			nodeMetrics, errors, nodeFields(), duration)

		go dockerCollector(m.ctx, node, fmt.Sprintf("http://%s:%d", node, dockerMetricsPort),
			dockerMetrics, errors, filter, duration)
	}

	if m.metrics == nil {
		return fmt.Errorf("metrics not set")
	}

	for {
		select {
		case nm := <-nodeMetrics:
			if stats, ok := m.nodeMetrics[nm.name]; ok {
				stats.update(CreateNodeStats(nm.name, nm.values))
			} else {
				nStats := CreateNodeStats(nm.name, nm.values)
				m.nodeMetrics[nm.name] = &nStats
			}
			m.metrics.Collect(m.nodeMetrics[nm.name].asMap())
		case dm := <-dockerMetrics:
			m.metrics.Collect(map[string]interface{}{
				"timestamp":      time.Now().Unix(),
				"HId":            dm.name,
				"docker_running": dm.Running,
				"docker_paused":  dm.Paused,
				"docker_total":   dm.Total,
			})
		case err := <-errors:
			fmt.Printf("%+v\n", err)
		case <-m.ctx.Done():
			return nil
		}
	}

}

func New(ctx context.Context, nodes []string) *Metrics {
	m := &Metrics{
		ctx:         ctx,
		nodes:       nodes,
		nodeMetrics: make(map[string]*NodeStats),
	}

	return m
}
