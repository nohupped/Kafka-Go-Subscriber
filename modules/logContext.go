package modules

import (
	"github.com/Sirupsen/logrus"
	"runtime"
)
// Function signature to modules.LogContext. Use this to log line numbers for debugging.
type GetContext func(*logrus.Entry) *logrus.Entry

// var Lc is similar to a function pointer to GetContext LogContext(logger *logrus.Entry) *logrus.Entry. Whenever a *logrus.Entry
// object is passed to this function, it will append extra fields file=<path to filename> func=<function name> line=<line no>
// This is good for debugging.
var Lc GetContext

// var Lnc is just a pointer to logrus.Entry object.
var Lnc *logrus.Entry

var Logger *logrus.Logger


// LogContext will use runtime.Caller and runtime.FuncForPC to enable filename, function name and linenumber in the log.
func LogContext(logger *logrus.Entry) *logrus.Entry {
	if pc, file, line, ok := runtime.Caller(1); ok {
		fName := runtime.FuncForPC(pc).Name()
		return logger.WithField("file", file).WithField("line", line).WithField("func", fName)
	} else {
		return logger
	}
}
