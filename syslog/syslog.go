// Package syslog wraps log/syslog package and logging avaiable globally.
package syslog

import (
	"log/syslog"
)

var log *syslog.Writer

// SetLogger attemps to connect to syslog with the given parameters
func SetLogger(network string, host string, port string,
	priority syslog.Priority, tag string) error {

	var err error
	log, err = syslog.Dial(network, host+":"+port, priority, tag)

	return err
}

// Alert logs an alert message to syslog
func Alert(m string) error {
	return log.Alert(m)
}

// Crit logs a critical message to syslog
func Crit(m string) error {
	return log.Crit(m)
}

// Debug logs a debug message to syslog
func Debug(m string) error {
	return log.Debug(m)
}

// Emerg logs an emergency message to syslog
func Emerg(m string) error {
	return log.Emerg(m)
}

// Err logs an error message to syslog
func Err(m string) error {
	return log.Err(m)
}

// Info logs an info message to syslog
func Info(m string) error {
	return log.Info(m)
}

// Notice logs a notice message to syslog
func Notice(m string) error {
	return log.Notice(m)
}

// Warning logs a warning message to syslog
func Warning(m string) error {
	return log.Warning(m)
}

// Close closes the connection to syslog
func Close() error {
	return log.Close()
}

func init() {
	SetLogger("tcp", "localhost", "514",
		syslog.LOG_WARNING|syslog.LOG_LOCAL0, "")
}
