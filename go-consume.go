package main

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
	"go-consume/modules"
	"flag"
)



var signals chan os.Signal

func init() {
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)


}

func main() {
	defer modules.LogFile.Close()
	configFile := flag.String("configfile", "/etc/go-sub.conf", "Path to the config file")
	flag.Parse()
	parsedConfig := modules.ParseConfig(*configFile)
	modules.InitLog(parsedConfig)
	modules.Logger.Infoln("Starting daemon")
	modules.Logger.Infoln("Loaded config files with the following values", parsedConfig.String())


	go func() {
		for ; ;  {
			select {
			case <-signals:
				modules.Logger.Warnln("Interrupt received, dying...")
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


	if err != nil {
		panic(err)
	}

	// Todo : declare msg outside loop and
	for {
		select {
		case msg := <-consumer.Messages():
			fmt.Fprintf(os.Stdout, "%s/%d/%d\n\t%s\n\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			consumer.MarkOffset(msg, "") // mark message as processed
		case errormessage := <- consumer.Errors():
			modules.Logger.Errorln(errormessage)
		case notificationmessage := <- consumer.Notifications():
			modules.Logger.Infoln("Rebalancing", notificationmessage.Claimed, notificationmessage.Current, notificationmessage.Released)
		}
	}
}
