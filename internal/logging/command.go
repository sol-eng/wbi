package logging

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	cmdlog *logrus.Logger
)

type MyFormatter struct {
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

func init() {
	// Setup the command output file
	logFile := "wbi-command-" + time.Now().Format("20060102T150405") + ".sh"
	f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	cmdlog = logrus.New()

	cmdlog.SetOutput(f)
	cmdlog.SetFormatter(new(MyFormatter))

	cmdlog.Info("#!/bin/bash\n")
}

// Info ...
func Info(format string, v ...interface{}) {
	cmdlog.Infof(format, v...)
}

// Warn ...
func Warn(format string, v ...interface{}) {
	cmdlog.Warnf(format, v...)
}

// Error ...
func Error(format string, v ...interface{}) {
	cmdlog.Errorf(format, v...)
}
