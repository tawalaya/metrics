# Measurement Utility
Small GoLang util to collect metrics form Prometheus and Docker

## Usage
Example:`metrics --nodes=systemA,systemB --name=foo --timeout=60s --interval=500ms` 
Collects metrics from systemA and systemB every 500ms for 60s and saves them to foo_YEAR_MONTH_DAY.csv

By default, each select node must provide access to the docker api over port 2376 and access to the prometheus node agent over port 9100.

| Name     | Default | Usage                                                        |
|----------|---------|--------------------------------------------------------------|
| nodes    | `nil`   | list of urls to connect to seperated with `,`                | 
| name     | ` `     | prefix used for the log file                                 |
| timeout  | `-1`    | timeout before stopping to collect data, -1 means no timeout |
| interval | `1s`    | interval between polling data from all nodes                 |
 | docker   | `true`  | indecate if docker metrics should be collected               |
 | dport    | `2376`  | port used to connect to docker api                           |
| pport    | `9100`  | port used to connect to prometheus node agent                |


Use `Ctr+C`  to terminate the process if no timeout is defined.

## Output 
The output is a csv file with the following columns:

| Name             | Description                           |
|------------------|---------------------------------------|
| timestamp        | current timestamp                     |
| HId              | node name                             |
| cpu_user         | cpu used  py user in percent          |
| cpu_system       | cpu used  py system in percent        |
| cpu_idle         | cpu idle  in percent                  |
| cpu_iowait       | cpu used  py iowait in percent        |
| forks            | number of  forks                      |
| disk             | disc usage                            |
| load             | running load  average (1 minute)      |
| mem_available    | available memory  in bytes            |
| mem_total        | total memory  in bytes                |
| net_rev          | bytes received                        |
| net_tra          | bytes transmitted                     |
| cpu_user_delta   | cpu user  changed since last update   |
| cpu_system_delta | cpu system  changed since last update |
| cpu_idle_delta   | cpu idle  changed since last update   |
| cpu_iowait_delta | cpu iowait  changed since last update |
| forks_delta      | forks changed  since last update      |
| mem_usage        | memory usage  in bytes                |
| docker_running   | number of  running docker containers  |
| docker_paused    | number of  paused docker containers   |
| docker_total     | total number  of docker containers    |
