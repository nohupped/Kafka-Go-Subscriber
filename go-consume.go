package main

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
//	"fmt"
	"github.com/Shopify/sarama"
	"time"
	"go-consume/modules"
	"flag"
	"go-consume/plugins"
)



var signals chan os.Signal

func init() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)


}

func main() {
	defer modules.LogFile.Close()
	configFile := flag.String("configfile", "/etc/go-sub.ini", "Path to the config file")
	flag.Parse()
	parsedConfig := modules.ParseConfig(*configFile)
	modules.InitLog(parsedConfig)
	modules.Logger.Infoln("Starting Service.....")
	modules.Logger.Infoln("Loaded config files with the following values", parsedConfig.String())

	MessagesSyslogChan := make(chan *sarama.ConsumerMessage, parsedConfig.Daemon.MessageBuffer)

	go func() {
		for ; ;  {
			select {
			case <-signals:
				modules.Logger.Warnln("Interrupt received, closing open channels and dying...")
				close(MessagesSyslogChan)
				os.Exit(1)
			}
		}
	}()

	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Millisecond
	config.Group.Return.Notifications = true
	consumer, err := cluster.NewConsumer(parsedConfig.Kafka.Brokers, parsedConfig.Kafka.Groupid,
		parsedConfig.Kafka.Topics, config)

	// Start plugin syslog. Whatever topics that needs to be processed by this plugin must write to MessagesSyslogChan.
	go func() {
	plugins.PluginSyslog(MessagesSyslogChan, consumer)
	}()


	if err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-consumer.Messages():
			switch msg.Topic {
			case "rsyslog":
				MessagesSyslogChan <- msg
				break
			}
		//	consumer.MarkOffset(msg, "") // mark message as processed
		case errormessage := <- consumer.Errors():
			modules.Logger.Errorln(errormessage)
		case notificationmessage := <- consumer.Notifications():
			modules.Logger.Infoln("Rebalancing", notificationmessage.Claimed, notificationmessage.Current, notificationmessage.Released)
		}
	}
}
