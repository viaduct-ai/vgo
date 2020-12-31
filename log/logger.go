package log

// Logger is a generic logging interface.
// It is motivated by
// https://dave.cheney.net/2015/11/05/lets-talk-about-logging
// https://github.com/golang/go/issues/13182
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	With(args ...interface{}) Logger
}
