
Sample config file:


```[Daemon]
SaveState = true

; Available log levels are 2(Error), 3(Warning, which includes Error and Warning),
; 4(Info, which includes Error, Warn and Info), 5(Debug, which includes Error, Warn, Info and Debug)
; Defaults to 2
loglevel = 4


[Kafka]
Brokers = "172.19.85.78:9092,172.19.40.17:9092,172.19.1.188:9092"
;topics = rsyslog, graphite-AutoOpt
topics = rsyslog ;graphite-AutoOpt
groupid = testrsyslog
```