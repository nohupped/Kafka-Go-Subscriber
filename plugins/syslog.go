package plugins

import (
	"github.com/Shopify/sarama"
	//"fmt"
	//"os"
	"github.com/bsm/sarama-cluster"
	log "github.com/nohupped/glog"
//	"github.com/jeromer/syslogparser"
//	"github.com/jeromer/syslogparser/rfc3164"
//	"github.com/jeromer/syslogparser/rfc5424"


	"time"
//	"fmt"
	//"os"
	"encoding/json"
	"github.com/jeromer/syslogparser/rfc5424"
	"fmt"
)

func PluginSyslog(messages chan *sarama.ConsumerMessage, consumer *cluster.Consumer, logger *log.Logger) {
	logger.Infoln("Started plugin syslog...")
	// Ranging over the channel messages.
	for msg := range messages {
		//fmt.Fprintf(os.Stdout, "From Plugin, %s/%d/%d\n\t%s\n\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		ParseSyslog(msg.Value)
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

func ParseSyslog(msg []byte) {
	//fmt.Fprintf(os.Stdout, "%s\n", msg)
	data := make(map[string]interface{})
	err := json.Unmarshal(msg, &data)
	if err != nil {
		panic(err)
	}
	strippedMsg := []byte(data["message"].(string))
	p1 := rfc5424.NewParser(strippedMsg)
	p1.Parse()
	fmt.Println(p1.Dump())
	/*	p1 := rfc5424.NewParser(msg)
	p2 := rfc3164.NewParser(msg)
	p1.Parse()
	p2.Parse()
	fmt.Println("From rfc5424", p1.Dump())
	fmt.Println("From rfc3164", p2.Dump())
	*/
}
