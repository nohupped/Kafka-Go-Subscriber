
Sample config file:

```
[Daemon]

; Available log levels are 0(Error), 1(Warning, which includes Error and Warning),
; 2(Info, which includes Error, Warn and Info), 3(Debug, which includes Error, Warn, Info and Debug)
; Defaults to 2
loglevel = 2
; Logfile to where the daemon should log. Defaults to /var/log/kafka-go-sub.log.
logfile = /var/log/kafka-go-sub.log
MessageBuffer = 2048

[plugins]
; Add enabled plugins in coma separated list
enabledplugins = syslog

[pluginmap]
; Provide a json mapping.
syslog = { "rsyslog": "syslog", "someothertopic": "Dummy" }


[Kafka]
Brokers = "broker1:9092,broker2:9092,broker3:9092"
; topics = rsyslog, someothertopic
topics = rsyslog
;  Change groupid to a new one to consume from starting again. The current 
; library doesn't provide a way to reset the offset to a lower value 
; than the current offset.
groupid = testnewplugin12
; plugin config

[syslog]
syslog_send_protocol = tcp
syslog_server_ip_port = "rsyslogserver:514"
syslog_server_dialtimeout = 10 ; in seconds
enable_offset_logging = true ; disable for minor performance improvements
offset_Logging_interval = 1 ; in minutes



```
