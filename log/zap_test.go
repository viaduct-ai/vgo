package log_test

import (
	"testing"

	"github.com/viaduct-ai/vgo/log"
)

type testLogger struct {
	log.Logger
	called map[string]bool
}

func (l *testLogger) Error(args ...interface{}) {
	l.called["Error"] = true
	l.Logger.Error(args...)
}

func (l *testLogger) Info(args ...interface{}) {
	l.called["Info"] = true
	l.Logger.Info(args...)
}

func (l *testLogger) Debug(args ...interface{}) {
	l.called["Debug"] = true
	l.Logger.Debug(args...)
}

func (l *testLogger) Errorf(template string, args ...interface{}) {
	l.called["Errorf"] = true
	l.Logger.Errorf(template, args...)
}

func (l *testLogger) Infof(template string, args ...interface{}) {
	l.called["Infof"] = true
	l.Logger.Infof(template, args...)
}

func (l *testLogger) Debugf(template string, args ...interface{}) {
	l.called["Debugf"] = true
	l.Logger.Debugf(template, args...)
}

func (l *testLogger) With(args ...interface{}) log.Logger {
	l.called["With"] = true
	l.Logger = l.Logger.With(args...)
	return l
}

func newTestLogger() testLogger {
	return testLogger{
		Logger: log.NewNoOpZapLogger(),
		called: map[string]bool{},
	}
}

func TestZapLogger(t *testing.T) {
	t.Parallel()
	logger := newTestLogger()

	args := "test"
	tests := []struct {
		name string
		call func()
	}{
		{
			name: "Debug",
			call: func() {
				logger.Debug(args)
			},
		},
		{
			name: "Debugf",
			call: func() {
				logger.Debugf("%s", args)
			},
		},
		{
			name: "Info",
			call: func() {
				logger.Info(args)
			},
		},
		{
			name: "Infof",
			call: func() {
				logger.Infof("%s", args)
			},
		},
		{
			name: "Error",
			call: func() {
				logger.Error(args)
			},
		},
		{
			name: "Errorf",
			call: func() {
				logger.Errorf("%s", args)
			},
		},
		{
			name: "With",
			call: func() {
				logger.With(args, args)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.call()

			if !logger.called[tt.name] {
				t.Errorf("expected %q to be called", tt.name)
			}
		})
	}
}
