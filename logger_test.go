package glog

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
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
		name: "InfoWithDepth",
		loggingFunc: func(l *Logger) {
			l.InfoWithDepth(1, "hello")
		},
	},
	{
		name: "InfolnWithDepth",
		loggingFunc: func(l *Logger) {
			l.InfolnWithDepth(1, "hello")
		},
	},
	{
		name: "InfofWithDepth",
		loggingFunc: func(l *Logger) {
			l.InfofWithDepth(1, "hello: %s", "<NAME>")
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
		name: "WarningWithDepth",
		loggingFunc: func(l *Logger) {
			l.WarningWithDepth(1, "hello")
		},
	},
	{
		name: "WarninglnWithDepth",
		loggingFunc: func(l *Logger) {
			l.WarninglnWithDepth(1, "hello")
		},
	},
	{
		name: "WarningfWithDepth",
		loggingFunc: func(l *Logger) {
			l.WarningfWithDepth(1, "hello: %s", "<NAME>")
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
	{
		name: "ErrorWithDepth",
		loggingFunc: func(l *Logger) {
			l.ErrorWithDepth(1, "hello")
		},
	},
	{
		name: "ErrorlnWithDepth",
		loggingFunc: func(l *Logger) {
			l.ErrorlnWithDepth(1, "hello")
		},
	},
	{
		name: "ErrorfWithDepth",
		loggingFunc: func(l *Logger) {
			l.ErrorfWithDepth(1, "hello: %s", "<NAME>")
		},
	},
}

func TestLogData(t *testing.T) {
	defer resetOutput(setBuffer())
	clearBackends()

	comm := RegisterBackend()
	done := make(chan struct{})
	go func() {
		for {
			select {
			case e, open := <-comm:
				if !open {
					return
				}
				if !strings.Contains(fmt.Sprintf("%v", e), "content to ignore") {
					t.Error("backend did not received expected data")
				}
			case <-done:
				return
			}
		}
	}()

	logger := WithData("content to ignore")
	logger.Error("interesting content")

	logger = NewLogger()
	logger = logger.WithData("content to ignore")
	logger.Error("extra content")

	if contains("content to ignore", t) {
		t.Error("glog did not ignore data which it was told to ignore")
	}
	if !contains("interesting content", t) {
		t.Error("glog ignored content it was not supposed to")
	}
	if !contains("extra content", t) {
		t.Error("glog ignored content it was not supposed to")
	}

	close(done)
}
