package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func setLogLevel(level string) {
	// We want the log to be reset whenever it is initialized.
	logLevel := strings.ToLower(level)

	switch logLevel {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.Fatalf("Invalid log level: %s", logLevel)
	}
}
