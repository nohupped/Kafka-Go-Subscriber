package modules

import "fmt"

type Config struct {
	// Kafka section from config file.
	Kafka kafka
	// Daemon section from config file.
	Daemon daemon
}



type kafka struct {
	Brokers []string
	Topics []string
	Groupid string
}
type daemon struct {
	Loglevel int
	Logfile string
}



func (c Config)String()  string{
	return fmt.Sprintf("%#v\n", c)
}
