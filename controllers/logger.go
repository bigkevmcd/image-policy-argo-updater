package controllers

import "github.com/go-logr/logr"

const debug = 1

type logger struct {
	log logr.Logger
}

func newLogger(l logr.Logger) logger {
	return logger{log: l}
}

func (l logger) debug(msg string, keysAndValues ...interface{}) {
	l.log.V(debug).Info(msg, keysAndValues...)
}

func (l logger) error(err error, msg string, keysAndValues ...interface{}) {
	l.log.Error(err, msg, keysAndValues...)
}

func (l logger) info(msg string, keysAndValues ...interface{}) {
	l.log.Info(msg, keysAndValues...)
}
