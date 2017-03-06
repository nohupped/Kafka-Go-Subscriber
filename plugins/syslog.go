package plugins

import (
	"github.com/Shopify/sarama"
	"fmt"
	"os"
	"github.com/bsm/sarama-cluster"
)

func PluginSyslog(messages chan *sarama.ConsumerMessage, consumer *cluster.Consumer) {
	for {
		select {
		case msg := <- messages:
			fmt.Fprintf(os.Stdout, "From Plugin, %s/%d/%d\n\t%s\n\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
			consumer.MarkOffset(msg, "")
		}
	}
}
