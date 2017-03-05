
Sample config file:

```
[Daemon]

; Available log levels are 0(Error), 1(Warning, which includes Error and Warning),
; 2(Info, which includes Error, Warn and Info), 3(Debug, which includes Error, Warn, Info and Debug)
; Defaults to 2
loglevel = 3
; Logfile to where the daemon should log. Defaults to /var/log/kafka-go-sub.log.
logfile = /var/log/kafka-go-sub.log


[Kafka]
Brokers = "172.19.85.78:9092,172.19.40.17:9092,172.19.1.188:9092"
topics = rsyslog, graphite-AutoOpt
groupid = testrsyslog
```