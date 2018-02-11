package glog

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestPrefix(t *testing.T) {
	for _, test := range loggerTests {
		t.Run(test.name, func(t *testing.T) {
			defer resetOutput(setBuffer())
			prefixedLogger := WithPrefix("examplePrefix")
			test.loggingFunc(prefixedLogger)
			if !contains("examplePrefix", t) {
				t.Errorf("Prefix was not included: %s", contents())
			}
		})
	}
}
func TestLoggerWithPrefix(t *testing.T) {
	for _, test := range loggerTests {
		t.Run(test.name, func(t *testing.T) {
			defer resetOutput(setBuffer())
			logger := NewLogger()
			prefixedLogger := logger.WithPrefix("examplePrefix")
			test.loggingFunc(prefixedLogger)
			if !contains("examplePrefix", t) {
				t.Errorf("Prefix was not included: %s", contents())
			}
		})
	}
}

func TestContextPrefix(t *testing.T) {
	for _, test := range loggerTests {
		t.Run(test.name, func(t *testing.T) {
			defer resetOutput(setBuffer())
			ctx := context.Background()
			ctx = ContextWithPrefix(ctx, "examplePrefix")
			prefixedLogger := WithContext(ctx)
			test.loggingFunc(prefixedLogger)
			if !contains("examplePrefix", t) {
				t.Errorf("Prefix was not included: %s", contents())
			}
		})
	}
}

var loggerTests = []struct {
	name        string
	loggingFunc func(l *Logger)
}{
	{
		name: "Info",
		loggingFunc: func(l *Logger) {
			l.Info("hello")
		},
	},
	{
		name: "Infoln",
		loggingFunc: func(l *Logger) {
			l.Infoln("hello")
		},
	},
	{
		name: "Infof",
		loggingFunc: func(l *Logger) {
			l.Infof("hello: %s", "<NAME>")
		},
	},
	{
		name: "Warning",
		loggingFunc: func(l *Logger) {
			l.Warning("hello")
		},
	},
	{
		name: "Warningln",
		loggingFunc: func(l *Logger) {
			l.Warningln("hello")
		},
	},
	{
		name: "Warningf",
		loggingFunc: func(l *Logger) {
			l.Warningf("hello: %s", "<NAME>")
		},
	},
	{
		name: "Error",
		loggingFunc: func(l *Logger) {
			l.Error("hello")
		},
	},
	{
		name: "ErrorIf",
		loggingFunc: func(l *Logger) {
			l.ErrorIf(errors.New("error"))
		},
	},
	{
		name: "Errorln",
		loggingFunc: func(l *Logger) {
			l.Errorln("hello")
		},
	},
	{
		name: "Errorf",
		loggingFunc: func(l *Logger) {
			l.Errorf("hello: %s", "<NAME>")
		},
	},
	{
		name: "ErrorfIf",
		loggingFunc: func(l *Logger) {
			l.ErrorfIf(errors.New("error"), "hello: %s", "<NAME>")
		},
	},
}

func TestLogData(t *testing.T) {
	defer resetOutput(setBuffer())

	comm := RegisterBackend()
	message1 := fmt.Sprintf("testLogData message: %v", time.Now().Nanosecond())

	logger := WithData("data1")
	logger.Error(message1)

	waitForData(t, comm, message1, "data1")

	message2 := fmt.Sprintf("testLogData message2: %v", time.Now().Nanosecond())

	logger = NewLogger()
	logger = logger.WithData("data2")
	logger.Error(message2)

	if contains("data1", t) || contains("data2", t) {
		t.Error("glog did not ignore data which it was told to ignore")
	}
	if !contains(message1, t) || !contains(message2, t) {
		t.Error("glog ignored content it was not supposed to")
	}

	waitForData(t, comm, message2, "data2")
}

func TestAppendData(t *testing.T) {
	defer resetOutput(setBuffer())

	comm := RegisterBackend()

	logger := WithData("data1")
	logger = logger.AppendData("data2")
	message := fmt.Sprintf("testAppendData message: %v", time.Now().Nanosecond())
	logger.Error(message)

	if contains("data1", t) || contains("data2", t) {
		t.Error("glog did not ignore data which it was told to ignore")
	}

	if !contains(message, t) {
		t.Error("glog ignored content it was not supposed to")
	}

	waitForData(t, comm, message, "data1", "data2")
}
