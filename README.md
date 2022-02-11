# Mesurment Utility

Small GoLang util to collect metrics form Prometheus and Docker

## Usage

`metrics --nodes=<node-ip/url>,<node-ip,url> --name <logname-prefix> --timeout=10s --interval=500ms`

Each select node must provide access to the docker socket over port 2376 and access to the prometheus node agent over port 9100.

| Name | Default | Usage |
| --- | ---- | --- |
| nodes | `nil` | list of urls to connect to seperated with `,` | 
| name | ` ` | prefix used for the log file |
| timeout | `-1` | timeout before stopping to collect data, -1 means no timeout |
| interval | `1s` | interval between polling data from all nodes|

Use `Ctr+C`  to terminate the process if no timeout is defined.
