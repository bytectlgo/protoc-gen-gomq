package main

import (
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
)

type loggerDebug struct {
	log *log.Helper
}

func (l *loggerDebug) Println(v ...interface{}) {
	l.log.Debug(v...)
}

func (l *loggerDebug) Printf(format string, v ...interface{}) {
	l.log.Debugf(format, v...)
}

type loggerError struct {
	log *log.Helper
}

func (l *loggerError) Println(v ...interface{}) {
	l.log.Error(v...)
}

func (l *loggerError) Printf(format string, v ...interface{}) {
	l.log.Errorf(format, v...)
}

type loggerWarn struct {
	log *log.Helper
}

func (l *loggerWarn) Println(v ...interface{}) {
	l.log.Warn(v...)
}

func (l *loggerWarn) Printf(format string, v ...interface{}) {
	l.log.Warnf(format, v...)
}

type loggerCritical struct {
	log *log.Helper
}

func (l *loggerCritical) Println(v ...interface{}) {
	l.log.Error(v...)
}

func (l *loggerCritical) Printf(format string, v ...interface{}) {
	l.log.Errorf(format, v...)
}

func init() {
	helper := log.NewHelper(log.With(log.DefaultLogger, "module", "pmqtt"))
	logDebug := &loggerDebug{
		log: helper,
	}
	logError := &loggerError{
		log: helper,
	}
	logWarn := &loggerWarn{
		log: helper,
	}
	logCritical := &loggerCritical{
		log: helper,
	}
	pmqtt.DEBUG = logDebug
	pmqtt.WARN = logWarn
	pmqtt.ERROR = logError
	pmqtt.CRITICAL = logCritical
}
