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
	"sync"
	"reflect"
)





func main() {
	defer modules.LogFile.Close()
	configFile := flag.String("configfile", "/etc/go-sub.ini", "Path to the config file")
	flag.Parse()
	parsedConfig := modules.ParseConfig(*configFile)
	modules.InitLog(parsedConfig)
	modules.Logger.Infoln("Starting Service.....")
	modules.Logger.Infoln("Loaded config files with the following values", parsedConfig.String())

	MessagesSyslogChan := make(chan *sarama.ConsumerMessage, parsedConfig.Daemon.MessageBuffer)
	wg := new(sync.WaitGroup)



	config := cluster.NewConfig()
	//config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Millisecond
	config.Group.Return.Notifications = true
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	consumer, err := cluster.NewConsumer(parsedConfig.Kafka.Brokers, parsedConfig.Kafka.Groupid,
		parsedConfig.Kafka.Topics, config)
	modules.Err(err)


	// Start plugins.
	for _, i := range parsedConfig.EnablePlugin.PluginsEnabled {
		switch i {
		// will start syslog plugin in a separate goroutine. Using a *sync.Waitgroup to the plugin running closure
		// so that when a signal is received, the plugin waits until the channel is closed, and all the messages
		// in the channel are processed.
		case "syslog":
			wg.Add(1)
			go func() {
				plugins.StartPluginSyslog(MessagesSyslogChan, consumer, parsedConfig.Syslog.SyslogSendProtocol,
				parsedConfig.Syslog.SyslogServerIPnPort, parsedConfig.Syslog.SyslogServerDialTimeout, modules.Logger,
				parsedConfig.Syslog.OffsetLogging, parsedConfig.Syslog.OffsetLoggingInterval)
				wg.Done()
			}()
			break
		default:

		}
	}




	for {
		select {
		case msg := <-consumer.Messages():
			switch parsedConfig.PluginMap.Syslog.TopicsToSyslog[msg.Topic] {
			// Additional plugin names to be evaluated here.
			case "syslog" :
				MessagesSyslogChan <- msg
				break
			}
		//	consumer.MarkOffset(msg, "") // mark message as processed
		case errormessage := <- consumer.Errors():
			modules.Logger.Errorln(errormessage)
		case notificationmessage := <- consumer.Notifications():
			modules.Logger.Infoln("Rebalancing", notificationmessage.Claimed, notificationmessage.Current, notificationmessage.Released)
		case <-signals:
		// Once a signal is received, the consumer is closed first, preventing any further reads from the kafka. Then the
		// channel is closed. Once the channel is closed, the syslog plugin's range over the channel will be finished
		// after processing all the messages in the channel's buffer and the function PluginSyslog will die.
		// Only after this, the sync.Waitgroup in the closure where the plugin is started will decrement,
		// causing the main thread to wait until the syslog plugin gracefully shutdown.
			modules.Logger.Infoln("Closing the consumer...")
			err := consumer.Close()
			modules.Logger.Infoln("Closed the consumer...")
			if err != nil {
				modules.Logger.Errorln(err)
			}
			modules.Logger.Infoln("Closing syslog channel...")
			close(MessagesSyslogChan)
			modules.Logger.Infoln("Closed syslog channel")
			wg.Wait()
			modules.Logger.Warnln("Interrupt received, Closed", reflect.TypeOf(MessagesSyslogChan),
				reflect.TypeOf(consumer), reflect.TypeOf(modules.LogFile))
			modules.Logger.Println("bye...")
		// Finally closing the log file handle
			err = modules.LogFile.Close()
			if err != nil {
				modules.Logger.Errorln(err)
			}
			os.Exit(1)

		}
	}
}
