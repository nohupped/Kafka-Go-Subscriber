
Sample config file:

```
[Daemon]

; Available log levels are 0(Error), 1(Warning, which includes Error and Warning),
; 2(Info, which includes Error, Warn and Info), 3(Debug, which includes Error, Warn, Info and Debug)
; Defaults to 2
loglevel = 3
; Logfile to where the daemon should log. Defaults to /var/log/kafka-go-sub.log.
logfile = /var/log/kafka-go-sub.log
MessageBuffer = 2048

[plugins]
; Add enabled plugins in coma separated list
enabledplugins = syslog

[pluginmap]
; Provide a json mapping.
syslog = { "rsyslog": "syslog", "graphite-AutoOpt": "Dummy" }


[Kafka]
Brokers = "172.19.85.78:9092,172.19.40.17:9092,172.19.1.188:9092"
; topics = rsyslog, graphite-AutoOpt
topics = rsyslog
;  Change groupid to a new one to consume from starting again. The current 
; library doesn't provide a way to reset the offset to a lower value 
; than the current offset.
groupid = testnewplugin3
; plugin config

[syslog]
syslog_send_protocol = tcp
syslog_server_ip_port = "172.24.40.36:514"
syslog_server_dialtimeout = 10 ; in seconds



```