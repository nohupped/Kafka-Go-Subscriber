package modules

import "fmt"

type Config struct {
	// Kafka section from config file.
	Kafka kafka
	// Daemon section from config file.
	Daemon daemon
	// PluginMaps holds a map of topics to its corresponding plugin to process the topics's messages.
	PluginMaps pluginMaps
}



type kafka struct {
	Brokers []string
	Topics []string
	Groupid string
}
type daemon struct {
	Loglevel int
	MessageBuffer int
	Logfile string

}
type pluginMaps struct {
	TopicsToPluginMap []map[string]string
}


func (c Config)String()  string{
	return fmt.Sprintf("%#v\n", c)
}

