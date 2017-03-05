package modules

import (
	log "github.com/nohupped/glog"
	"os"
)

var LogFile *os.File
var Logger *log.Logger
func InitLog(parsedConfig *Config) *log.Logger{
	var err error
	LogFile, err = os.OpenFile(parsedConfig.Daemon.Logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	Err(err)
	Logger = log.New(LogFile, "", log.Lshortfile|log.Ldate|log.Ltime)
	Logger.SetOutput(LogFile)
	Logger.Println("Setting loglevel to ", parsedConfig.Daemon.Loglevel)
	Logger.SetLogLevel(parsedConfig.Daemon.Loglevel)
	return Logger
}

