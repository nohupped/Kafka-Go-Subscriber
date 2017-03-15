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
func StartPluginSyslog(messages chan *sarama.ConsumerMessage, consumer *cluster.Consumer, syslogProto,
			syslogServernPort string, dialtimeout int, logger *log.Logger, enableOffsetLogging bool, offsetLoggingInterval ...int) {
	logger.Infoln("Started plugin syslog...")

//	var lastconsumedmessage []byte
//	var lastconsumedoffset int64
	// has to handle rsyslog daemon restarts during logrotates, which can cause a "connection reset by peer" error.
	DialSyslogServer(syslogProto, syslogServernPort, time.Second * time.Duration(dialtimeout), logger)

	// Ranging over the channel messages.

	switch enableOffsetLogging {
	case false:

		logger.Infoln("Offset logging is disabled, and will not be reported.")
		OutFalse:
		for msg := range messages {
			parsedMsg := ParseSyslog(msg.Value, logger)
			if parsedMsg == nil {
				logger.Errorln("Parsing returned nil for message", string(msg.Value), "at offset", msg.Offset, "for partition", msg.Partition, "on topic", msg.Topic, "Timestamp:", msg.Timestamp)
				continue
			}
			//parsedMsg := syslogConstructor.ConstructSyslogFromByteArray(msg.Value, "\",\"offset\"")
			var written int
			var WriteError error
			written, WriteError = conn.Write(parsedMsg)
			if WriteError != nil {
				logger.Errorln("Error", WriteError, "encountered when writing", string(parsedMsg), "at offset", msg.Offset,  "to socket, redialling...")
				DialSyslogServer(syslogProto, syslogServernPort, time.Second * time.Duration(dialtimeout), logger)
				written, WriteError = conn.Write(parsedMsg)
				if WriteError != nil {
					logger.Errorln("Cannot write message to socket, re-connection failed, plugin syslog dying...")
					break OutFalse
				}
				logger.Infoln("Resumed write of message", string(parsedMsg), "from offset", msg.Offset)
			}

			consumer.MarkOffset(msg, "")
//			lastconsumedmessage = msg.Value
//			lastconsumedoffset = msg.Offset
			logger.Debugln(written, "bytes written for the message", string(parsedMsg), "and marked consumer offset")
			//logger.Debugln(written, "bytes written for the message", parsedMsg, "and marked consumer offset")

		}


	case true:
		logger.Infoln("OffsetLogging is enabled and will be reported every", offsetLoggingInterval[0], "minutes once the consumer queue is read..")
		var offset int64
		offset = -1
		go func() {
			time.Sleep(time.Second * 5) // Adding a delay to logging
			var previousOffset int64
			previousOffset = -1

			for {
				if offset != -1 {
					logger.Infoln("Last read offset is", offset)
				}
				if previousOffset != -1 {
					logger.Infoln("Consumed", offset - previousOffset, "messages in", offsetLoggingInterval[0], "minutes.")
				}
				previousOffset = offset

				time.Sleep(time.Minute * time.Duration(offsetLoggingInterval[0]))
			}
		}()

		OutTrue:
		for msg := range messages {

			parsedMsg := ParseSyslog(msg.Value, logger)
			if parsedMsg == nil {
				logger.Errorln("Parsing returned nil for message", string(msg.Value), "at offset", msg.Offset, "for partition", msg.Partition, "on topic", msg.Topic, "Timestamp:", msg.Timestamp)
				continue
			}
			//parsedMsg := syslogConstructor.ConstructSyslogFromByteArray(msg.Value, "\",\"offset\"")
			var written int
			var WriteError error
			written, WriteError = conn.Write(parsedMsg)
			if WriteError != nil {
				logger.Errorln("Error", WriteError, "encountered when writing", string(parsedMsg), "at offset", msg.Offset,  "to socket, redialling...")
				DialSyslogServer(syslogProto, syslogServernPort, time.Second * time.Duration(dialtimeout), logger)
				written, WriteError = conn.Write(parsedMsg)
				if WriteError != nil {
					logger.Errorln("Error encountered", WriteError, "Cannot write message to socket, re-connection failed, plugin syslog dying...")
					break OutTrue
				}
				logger.Infoln("Resumed write of message", string(parsedMsg), "from offset", msg.Offset)

			}
			offset = msg.Offset

			consumer.MarkOffset(msg, "")
//			lastconsumedmessage = msg.Value
//			lastconsumedoffset = msg.Offset
			logger.Debugln(written, "bytes written for the message", string(parsedMsg), "and marked consumer offset")
			//logger.Debugln(written, "bytes written for the message", parsedMsg, "and marked consumer offset")
		}


	} // switch case closed
	// Once the consumer is closed upon receiving an interrupt, the range over the channel will be finished, and the below wg.Done
	// will decrement the waitgroup. This will make sure that the main
	logger.Infoln("Range over syslog channel exited because the channel was closed from elsewhere. Shutting down syslog plugin in 5 seconds ..")
//	logger.Infoln("Last consumed message:", string(lastconsumedmessage))
//	logger.Infoln("Last consumed offset:", lastconsumedoffset)
	logger.Infoln("Last committed offset may be of earlier offset than the consumed offset..")
	for i := 5; i >=0; i -- {
		logger.Infoln(i)
		time.Sleep(time.Second * 1)
	}
}
var conn net.Conn
//
func DialSyslogServer(proto, ipAndPort string, timeout time.Duration, logger *log.Logger) {
	logger.Infoln("Connecting to", ipAndPort)
	for i := 1; i <= 10; i++ {
		var DialError error
		conn, DialError = net.DialTimeout(proto, ipAndPort, timeout)
		if DialError != nil {
			logger.Errorln("Error when dialling to", ipAndPort, DialError, ", attempt", i)
			if i == 10 {
				logger.Errorln("All retries failed of error", DialError, ", dying..")
				panic(DialError)
			}
			time.Sleep(time.Second * 5)

		} else {
			logger.Infoln("Connection established...")
			break
		}
	}


}


func ParseSyslog(msg []byte, logger *log.Logger) []byte{
	buff := new(bytes.Buffer)
	buff.WriteString("<13>")
	data := make(map[string]interface{})
	err := json.Unmarshal(msg, &data)
	if err != nil {
		logger.Errorln("Error when processing", string(msg), err)
		return nil
	}
	buff.WriteString(data["message"].(string))
	buff.WriteString("\n")
	return buff.Bytes()


}
