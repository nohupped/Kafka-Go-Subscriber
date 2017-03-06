package modules

import (
	"gopkg.in/ini.v1"
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


/*	for _, i := range config.Kafka.Topics {

		pmap := make(map[string]string)
		plugin, err := pluginSection.GetKey(i)
		Err(err)

		pmap[i] = plugin.String()
		config.PluginMaps.TopicsToPluginMap = append(config.PluginMaps.TopicsToPluginMap, pmap)
	}*/


	return config

}