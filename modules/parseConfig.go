package modules

import (
	"gopkg.in/ini.v1"
	"fmt"
	"encoding/json"
)

func ParseConfig(path string) *Config {
	config := new(Config)
	loadOptions := ini.LoadOptions{AllowBooleanKeys: true, Insensitive: true}
	cfg, err := ini.LoadSources(loadOptions, path)
	Err(err)
	// Get section Kafka
	kafkaSection, _ := cfg.GetSection("kafka")
	brokers, err := kafkaSection.GetKey("brokers")
	Err(err)
	config.Kafka.Brokers = brokers.Strings(",")
	topics, err := kafkaSection.GetKey("topics")
	Err(err)
	config.Kafka.Topics = topics.Strings(",")
	groupid, err := kafkaSection.GetKey("groupid")
	Err(err)
	config.Kafka.Groupid = groupid.String()


	// Get section Daemon
	daemonSection, _ := cfg.GetSection("daemon")
	loglevel, err := daemonSection.GetKey("loglevel")
	Err(err)
	config.Daemon.Loglevel = loglevel.MustInt(2)
	logfile, err := daemonSection.GetKey("logfile")
	Err(err)
	config.Daemon.Logfile = logfile.MustString("/var/log/kafka-go-sub.log")
	messageBuffer, err := daemonSection.GetKey("MessageBuffer")
	config.Daemon.MessageBuffer = messageBuffer.MustInt(1024)
	// Generate plugin configuration for each topic
	pluginSection, err := cfg.GetSection("plugins")
	Err(err)

	enabledPlugins, err := pluginSection.GetKey("enabledplugins")
	config.EnablePlugin.PluginsEnabled = enabledPlugins.Strings(",")

	pluginMapSection, err := cfg.GetSection("pluginmap")
	b := make(map[string]string)
	for _, i := range pluginMapSection.Keys() {
		fmt.Println(i.String())
		err := json.Unmarshal([]byte(i.String()), &b)
		Err(err)
		config.PluginMap.Syslog.TopicsToSyslog = b
	}


	// Plugin configs

	syslogPluginSection, _ := cfg.GetSection("syslog")
	syslogProtocol, err := syslogPluginSection.GetKey("syslog_send_protocol")
	Err(err)
	config.Syslog.SyslogSendProtocol = syslogProtocol.MustString("tcp")
	syslogServerIPnPort, err := syslogPluginSection.GetKey("syslog_server_ip_port")
	Err(err)
	config.Syslog.SyslogServerIPnPort = syslogServerIPnPort.MustString("127.0.0.1:514")
	syslogServerDialTimeout, err := syslogPluginSection.GetKey("syslog_server_dialtimeout")
	config.Syslog.SyslogServerDialTimeout = syslogServerDialTimeout.MustInt(20)
	enableOffsetLogging, err := syslogPluginSection.GetKey("enable_offset_logging")
	Err(err)
	config.Syslog.OffsetLogging = enableOffsetLogging.MustBool(false)
	offsetLoggingTime, err := syslogPluginSection.GetKey("offset_Logging_Interval")
	Err(err)
	config.Syslog.OffsetLoggingInterval = offsetLoggingTime.MustInt(1)

	return config

}