# mackerel-plugin-postfix-log

Read and analyze postfix logs

## Usage

```
Usage:
  mackerel-plugin-postfix-log [OPTIONS]

Application Options:
      --logfile=        path to nginx ltsv logfile (/var/log/maillog)
      --posfile-prefix= prefix added position file (maillog)
  -v, --version         Show version

Help Options:
  -h, --help            Show this help message
 ```

Sample

```
postfixlog.total_delay.average  0.240476        1555681849
postfixlog.total_delay.99_percentile    0.250000        1555681849
postfixlog.total_delay.95_percentile    0.240000        1555681849
postfixlog.total_delay.90_percentile    0.240000        1555681849
postfixlog.recving_delay.average        0.040000        1555681849
postfixlog.recving_delay.99_percentile  0.040000        1555681849
postfixlog.recving_delay.95_percentile  0.040000        1555681849
postfixlog.recving_delay.90_percentile  0.040000        1555681849
postfixlog.queuing_delay.average        0.000476        1555681849
postfixlog.queuing_delay.99_percentile  0.010000        1555681849
postfixlog.queuing_delay.95_percentile  0.000000        1555681849
postfixlog.queuing_delay.90_percentile  0.000000        1555681849
postfixlog.connection_delay.average     0.090000        1555681849
postfixlog.connection_delay.99_percentile       0.090000        1555681849
postfixlog.connection_delay.95_percentile       0.090000        1555681849
postfixlog.connection_delay.90_percentile       0.090000        1555681849
postfixlog.transmission_delay.average   0.090000        1555681849
postfixlog.transmission_delay.99_percentile     0.090000        1555681849
postfixlog.transmission_delay.95_percentile     0.090000        1555681849
postfixlog.transmission_delay.90_percentile     0.090000        1555681849
postfixlog.transfer_num.2xx_count       1.615385        1555681849
postfixlog.transfer_num.4xx_count       0.000000        1555681849
postfixlog.transfer_num.5xx_count       0.000000        1555681849
postfixlog.transfer_total.count 1.615385        1555681849
postfixlog.transfer_ratio.2xx_percentage        100.000000      1555681849
postfixlog.transfer_ratio.4xx_percentage        0.000000        1555681849
postfixlog.transfer_ratio.5xx_percentage        0.000000        1555681849
```

## Install

Please download release page or `mkr plugin install monitoring-forge/mackerel-plugin-postfix-log`.