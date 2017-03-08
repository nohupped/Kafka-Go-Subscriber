package plugins

import (
	"github.com/Shopify/sarama"
	//"fmt"
	//"os"
	"github.com/bsm/sarama-cluster"
	log "github.com/nohupped/glog"

	"time"
	"net"
	"bytes"

)

// RFC syslog format: <13>2017-03-08T12:29:02.231335+00:00 openldap1 slapd[392]: slap_global_control: unrecognized control: 1.3.6.1.4.1.42.2.27.8.5.1
// Construct a log of similar format IF it is just a file tail from syslog or similar logs.
func StartPluginSyslog(messages chan *sarama.ConsumerMessage, consumer *cluster.Consumer, logger *log.Logger) {
	logger.Infoln("Started plugin syslog...")

	DialSyslogServer("tcp", "172.19.116.144:514", time.Second * 10)
	// Ranging over the channel messages.

	Out:
	for msg := range messages {
		//fmt.Fprintf(os.Stdout, "From Plugin, %s/%d/%d\n\t%s\n\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
		//strippedMsg := ParseSyslog(msg.Value)

		parsedMsg := ParseSyslog(msg.Value)
		var written int
		var WriteError error
		written, WriteError = conn.Write(parsedMsg)
		if WriteError != nil {
			DialSyslogServer("tcp", "172.19.116.144:514", time.Second * 10)
			written, WriteError = conn.Write(parsedMsg)
			if WriteError != nil {
				logger.Errorln("Cannor write message to socket, re-connection failed, plugin syslog dying...")
				break Out
			}
		}
		//if strippedMsg != nil {
		//	written, err := conn.Write(strippedMsg)

		//	if err != nil {
		//		panic(err)
		//	}
		//	fmt.Println("Written", written, "bytes")

		//}
//		fmt.Println(string(strippedMsg))
		consumer.MarkOffset(msg, "")
		logger.Debugln(written, "bytes written for the message", string(parsedMsg), "and marked consumer offset")

	}
	// Once the consumer is closed upon receiving an interrupt, the range over the channel will be finished, and the below wg.Done
	// will decrement the waitgroup. This will make sure that the main
	logger.Infoln("Range over syslog channel exited because the channel was closed from elsewhere. Shutting down syslog plugin in 5 seconds ..")
	for i := 5; i >=0; i -- {
		logger.Infoln(i)
		time.Sleep(time.Second * 1)
	}
}
var conn net.Conn
var DialError error
func DialSyslogServer(proto, ipAndPort string, timeout time.Duration)  {
	conn, DialError = net.DialTimeout(proto, ipAndPort, timeout)
	if DialError != nil {
		panic(DialError)
	}
}


func ParseSyslog(msg []byte) []byte{
	buff := new(bytes.Buffer)
	buff.WriteString("<13>")
	buff.Write(msg)
	return buff.Bytes()
//	data := make(map[string]interface{})
//	fmt.Println(string(msg))
	/*err := json.Unmarshal(msg, &data)
	if err != nil {
		panic(err)
	}
	//fmt.Println(data)
	buff.Write([]byte(data["message"].(string)))
	fmt.Println(buff.String())
	*/


}
