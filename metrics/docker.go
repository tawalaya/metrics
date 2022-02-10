package collector

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func dockerCollector(ctx context.Context, name, url string, stream chan DockerMetrics, errors chan error, filters filters.Args, refreshInterval time.Duration) {
	cli, err := client.NewClientWithOpts(client.WithHost(url), client.WithVersion("v1.24"))
	if err != nil {
		errors <- err
		return
	}
	for {
		list, err := cli.ContainerList(ctx, types.ContainerListOptions{
			Quiet:   false,
			Size:    false,
			All:     true,
			Latest:  false,
			Since:   "",
			Before:  "",
			Limit:   0,
			Filters: filters,
		})

		if err != nil {
			errors <- err
		} else {
			var running = 0
			var paused = 0
			var total = 0
			for _, container := range list {
				switch status := container.State; status {
				case "running":
					running++
				case "paused":
					paused++

				}
				total++
			}

			stream <- DockerMetrics{
				name:    name,
				Running: running,
				Paused:  paused,
				Total:   total,
			}
		}

		select {
		case _ = <-time.After(refreshInterval):
			//nothing
		case _ = <-ctx.Done():
			fmt.Printf("[%s] cacneld", name)
			return

		}
	}

}
