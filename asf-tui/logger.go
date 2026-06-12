package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var asfLog *log.Logger

func initLogger() error {
	logPath := asfLogPath()
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	multi := io.MultiWriter(f, os.Stderr)
	asfLog = log.New(multi, "[asf] ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}
