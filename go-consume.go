package main

import (
	"github.com/bsm/sarama-cluster"
	"os"
	"os/signal"
	"fmt"
	"github.com/Shopify/sarama"
	"time"
	"go-consume/modules"
	"github.com/Sirupsen/logrus"
	"flag"
)



var f *os.File
var signals chan os.Signal

func init() {
	//Initiate log
	var err error
	modules.Logger = logrus.New()
	f, err = os.OpenFile("/var/log/kafka-go-sub.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	modules.Err(err)
	modules.Logger.Out = f

	modules.Lc = modules.LogContext

	modules.Lnc = logrus.NewEntry(modules.Logger)
	modules.Lc(modules.Lnc).Infoln("Info from init")
	//modules.Lc.Infoln("Shit")

	logrus.SetOutput(f)
	signals = make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)


}

func main() {
	defer f.Close()
	modules.Lc(modules.Lnc).Infoln("With context")
	modules.Lnc.Infoln("Without context")
	//l(lentry).Infoln("Another test log")
	configFile := flag.String("configfile", "/etc/go-sub.conf", "Path to the config file")
	flag.Parse()
	parsedConfig := modules.ParseConfig(*configFile)
	modules.Lnc.Infoln(*parsedConfig)

	modules.Logger.Level = logrus.Level(parsedConfig.Daemon.Loglevel)
	modules.Lc(modules.Lnc).Debugln("Another log line")

	config := cluster.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.CommitInterval = 1 * time.Millisecond
	config.Group.Return.Notifications = true
	//brokers := []string{"172.19.85.78:9092","172.19.40.17:9092","172.19.1.188:9092"}
	//topics := []string{ "graphite-AutoOpt", "rsyslog"}
	consumer, err := cluster.NewConsumer(parsedConfig.Kafka.Brokers, parsedConfig.Kafka.Groupid,
		parsedConfig.Kafka.Topics, config)



	if err != nil {
		panic(err)
	}


	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	for {
		select {
		case msg := <-consumer.Messages():
			fmt.Fprintf(os.Stdout, "%s/%d/%d\n\t%s\n\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		//consumer.CommitOffsets()
//			consumer.MarkOffset(msg, "") // mark message as processed
			err := consumer.CommitOffsets() // mark message as processed
			if err != nil {
				fmt.Println(err)
			}
		case errormessage := <- consumer.Errors():
			fmt.Fprintf(os.Stderr, "%s\n", errormessage)
		case notificationmessage := <- consumer.Notifications():
			modules.Lnc.Infoln("Rebalancing", notificationmessage.Claimed, notificationmessage.Current, notificationmessage.Released)
		case <-signals:
			return
		}
	}
}
