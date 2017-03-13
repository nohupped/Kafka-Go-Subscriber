package plugins

import (
	"github.com/Shopify/sarama"

	"github.com/bsm/sarama-cluster"
	log "github.com/nohupped/glog"

	"time"
	"net"
	"bytes"

//	"syslogConstructor"
	"encoding/json"
)

// RFC syslog format: <13>2017-03-08T12:29:02.231335+00:00 openldap1 slapd[392]: slap_global_control: unrecognized control: 1.3.6.1.4.1.42.2.27.8.5.1
// Construct a log of similar format IF it is just a file tail from syslog or similar logs.
func StartPluginSyslog(messages chan *sarama.ConsumerMessage, consumer *cluster.Consumer, syslogProto, syslogServernPort string, dialtimeout int, logger *log.Logger) {
	logger.Infoln("Started plugin syslog...")

	DialSyslogServer(syslogProto, syslogServernPort, time.Second * time.Duration(dialtimeout))
	// Ranging over the channel messages.

	Out:
	for msg := range messages {


		parsedMsg := ParseSyslog(msg.Value)
		//parsedMsg := syslogConstructor.ConstructSyslogFromByteArray(msg.Value, "\",\"offset\"")
		var written int
		var WriteError error
		written, WriteError = conn.Write(parsedMsg)
		if WriteError != nil {
			DialSyslogServer(syslogProto, syslogServernPort, time.Second * time.Duration(dialtimeout))
			written, WriteError = conn.Write(parsedMsg)
			if WriteError != nil {
				logger.Errorln("Cannot write message to socket, re-connection failed, plugin syslog dying...")
				break Out
			}
		}

		consumer.MarkOffset(msg, "")
		logger.Debugln(written, "bytes written for the message", string(parsedMsg), "and marked consumer offset")
		//logger.Debugln(written, "bytes written for the message", parsedMsg, "and marked consumer offset")

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
	data := make(map[string]interface{})
	err := json.Unmarshal(msg, &data)
	if err != nil {
		panic(err)
	}
	buff.WriteString(data["message"].(string))
	buff.WriteString("\n")
	return buff.Bytes()


}
