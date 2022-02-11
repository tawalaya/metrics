package main

import (
	"context"
	"flag"
	"github.com/docker/docker/api/types/filters"
	"github.com/spf13/viper"
	collector "github.com/tawalaya/metrics/metrics"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var Nodes string
var name string
var timeout time.Duration = time.Duration(-1)
var interval time.Duration = time.Second

var dockerPort int
var prometheusPort int
var logDocker bool

func main() {
	flag.StringVar(&Nodes, "nodes", Nodes, "Nodes to connect to")

	flag.StringVar(&name, "name", "", "file name to write out or empty")
	flag.DurationVar(&timeout, "timeout", timeout, "timeout before stopping to collect data, -1 means no timeout")
	flag.DurationVar(&interval, "interval", interval, "interval between collecting data")

	flag.IntVar(&dockerPort, "dport", 2376, "docker port")
	flag.IntVar(&prometheusPort, "pport", 9100, "prometheus port")
	flag.BoolVar(&logDocker, "docker", true, "collect docker metrics")

	flag.Parse()

	if name != "" {
		viper.Set("logName", name)
	}

	if Nodes == "" {
		flag.PrintDefaults()
		return
	}

	workerNodes := strings.Split(Nodes, ",")

	var ctx context.Context
	var cancel context.CancelFunc

	if timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-s
		cancel()
		<-time.After(time.Second)
		os.Exit(1)
	}()

	c := collector.New(ctx, workerNodes)
	err := c.Setup(nil, logDocker, &dockerPort, &prometheusPort)
	if err != nil {
		panic(err)
	}

	err = c.Collect(interval, filters.Args{})
	if err != nil {
		log.Printf("failed to collect metrics: %v", err)
		return
	}
	defer cancel()

}
