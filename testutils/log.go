package testutils

import (
	"fmt"

	"github.com/viaduct-ai/vgo/log"
)

// TestLogger implements the log.Logger interface
// plus additional public fields for accessing what has been "written" to output.
// Debug, Info, and Error all have their own "output buffers".
type TestLogger struct {
	Context   map[string]interface{}
	DebugLogs []interface{}
	InfoLogs  []interface{}
	ErrorLogs []interface{}
}

// NewTestLogger returns a reference to new TestLogger
func NewTestLogger() *TestLogger {
	return &TestLogger{
		Context:   map[string]interface{}{},
		DebugLogs: []interface{}{},
		InfoLogs:  []interface{}{},
		ErrorLogs: []interface{}{},
	}
}

// Error logs args to the ErrorLogs slices
func (l *TestLogger) Error(args ...interface{}) {
	l.ErrorLogs = append(l.ErrorLogs, args...)
}

// Info logs args to the InfoLogs slices
func (l *TestLogger) Info(args ...interface{}) {
	l.InfoLogs = append(l.InfoLogs, args...)
}

// Debug logs args to the DebugLogs slices
func (l *TestLogger) Debug(args ...interface{}) {
	l.DebugLogs = append(l.DebugLogs, args...)
}

// Errorf logs the template string to the ErrorLogs slice
func (l *TestLogger) Errorf(template string, args ...interface{}) {
	str := fmt.Sprintf(template, args...)
	l.ErrorLogs = append(l.ErrorLogs, str)
}

// Infof logs the template string to the InfoLogs slice
func (l *TestLogger) Infof(template string, args ...interface{}) {
	str := fmt.Sprintf(template, args...)
	l.InfoLogs = append(l.InfoLogs, str)
}

// Debugf logs the template string to the DebugLogs slice
func (l *TestLogger) Debugf(template string, args ...interface{}) {
	str := fmt.Sprintf(template, args...)
	l.DebugLogs = append(l.DebugLogs, str)
}

// With mutates the existing TestLogger to include the additional context and returns it.
func (l *TestLogger) With(args ...interface{}) log.Logger {
	for i, v := range args {
		// even = key
		// odd = value
		if i%2 == 1 {
			// ignore ok b/c we checked above
			key, ok := args[i-1].(string)

			if !ok {
				panic("TestLogger received a non-string key")
			}
			l.Context[key] = v
		}
	}

	return l
}
