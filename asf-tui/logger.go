package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var asfLog *log.Logger
var debugLog *log.Logger

func initLogger() error {
	logPath := asfLogPath()
	if p := os.Getenv("ASF_LOG_FILE"); p != "" {
		logPath = p
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	var asfWriter io.Writer = f
	var debugWriter io.Writer = f

	debug := os.Getenv("ASF_DEBUG")
	if debug != "" && debug != "0" && debug != "false" && debug != "no" {
		asfWriter = io.MultiWriter(f, os.Stderr)
		debugWriter = io.MultiWriter(f, os.Stderr)
	}

	asfLog = log.New(asfWriter, "[asf] ", log.Ldate|log.Ltime|log.Lshortfile)
	debugLog = log.New(debugWriter, "[asf-debug] ", log.Ltime|log.Lshortfile)
	return nil
}
