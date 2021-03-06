package journal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	stdlog "log"

	kitlog "github.com/go-kit/kit/log"
)

var (
	logger  kitlog.Logger
	Service string
)

func init() {
	logger = kitlog.NewJSONLogger(os.Stdout)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))
}

// SetLogger allows user to output log data to a writer.
func SetLogger(w io.Writer) {
	logger = kitlog.NewJSONLogger(w)
	stdlog.SetOutput(kitlog.NewStdlibAdapter(logger))
}

// SetLogFile allows the user to output log data to a file.
func SetLogFile(file string) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		LogError(fmt.Sprintf("error opening log file: %v", err))
	}
	SetLogger(f)
}

// GetLogFile sets the log file, and tries to
// create the file if it doesn't exist.
func GetLogFile(logFile, fallback string) (string, error) {
	// The log file can be supplied by an
	// environment variable (for example),
	// so we also send a fallback in case that is ever empty.
	if logFile == "" && fallback == "" {
		return "", errors.New("please supply a log file name and fallback file")
	}

	if logFile == "" {
		logFile = fallback
	}

	if _, err := os.Stat(logFile); err != nil {
		LogError(fmt.Sprintf("log file %s does not exist", logFile))
		_, err := os.Create(logFile)
		if err != nil {
			return "", errors.New(fmt.Sprintf("unable to open log file: %s\n", logFile))
		}
		return logFile, nil
	}
	return logFile, nil
}

// LogRequest logs details of an HTTP request.
func LogRequest(r *http.Request) {
	logger.Log("channel", "request", "service", Service, "method", r.Method, "url", r.URL.String(), "headers", r.Header, "ts", time.Now())
}

// LogChannel logs data to a log channel.
func LogChannel(channel string, message ...interface{}) {
	logger.Log("channel", channel, "service", Service, "message", message, "ts", time.Now())
}

// LogError logs error data to the error channel
func LogError(message string) {
	LogChannel("error", message)
}

// LogInfo logs informational data.
func LogInfo(message string) {
	LogChannel("information", message)
}

// LogWorker logs information around queue worker operations.
func LogWorker(message ...interface{}) {
	LogChannel("worker", message)
}
