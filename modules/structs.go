package modules

import "fmt"

type Config struct {
	// Kafka section from config file.
	Kafka kafka
	// Daemon section from config file.
	Daemon daemon
	// EnablePlugins contain a list of plugins that are configured to run. A goroutine will spin up
	// for each enabled plugin.
	EnablePlugin enablePlugin

	// PluginMap has
	PluginMap pluginMap

}

type pluginMap struct {
	Syslog struct{
		TopicsToSyslog map[string]string
	       }
}

type enablePlugin struct {
	PluginsEnabled []string
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


func (c Config)String()  string {
	return fmt.Sprintf("%#v\n", c)
}


