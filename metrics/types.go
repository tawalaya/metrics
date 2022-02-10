package collector

import (
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type NodeMetrics struct {
	name   string
	values map[string]*io_prometheus_client.MetricFamily
}

type DockerMetrics struct {
	name    string
	Running int `json:"running"`
	Paused  int `json:"paused"`
	Total   int `json:"total"`
}
