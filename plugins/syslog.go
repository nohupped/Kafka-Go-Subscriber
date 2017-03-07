package plugins

import (
	"github.com/Shopify/sarama"
	"fmt"
	"os"
	"github.com/bsm/sarama-cluster"
	log "github.com/nohupped/glog"

	"time"
)

func PluginSyslog(messages chan *sarama.ConsumerMessage, consumer *cluster.Consumer, logger *log.Logger) {
	logger.Infoln("Started plugin syslog...")
	// Ranging over the channel messages.
	for msg := range messages {
		fmt.Fprintf(os.Stdout, "From Plugin, %s/%d/%d\n\t%s\n\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		consumer.MarkOffset(msg, "")
	}
	// Once the consumer is closed upon receiving an interrupt, the range over the channel will be finished, and the below wg.Done
	// will decrement the waitgroup. This will make sure that the main
	logger.Infoln("Range over syslog channel exited because the channel was closed from elsewhere. Shutting down syslog plugin in 5 seconds ..")
	for i := 5; i >=0; i -- {
		logger.Infoln(i)
		time.Sleep(time.Second * 1)
	}
}
