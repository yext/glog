package glog

import (
	"context"
	"fmt"
)

// Logger provides logging functionality with additional data and prefixing.
type Logger struct {
	*loggingT
	prefix string
	// Extra arguments to be appended to each request
	data []interface{}
}

// NewLogger creates a Logger instance with no additional data.
func NewLogger() *Logger {
	return &Logger{
		loggingT: &logging,
	}
}

// WithContext creates a logger from a context.Context
func WithContext(ctx context.Context) *Logger {
	data, _ := ctx.Value(contextKeyData).([]interface{})
	prefix, _ := ctx.Value(contextKeyPrefix).(string)
	return &Logger{
		loggingT: &logging,
		data:     data,
		prefix:   prefix,
	}
}

// WithPrefix creates a Logger with a given prefix.
func WithPrefix(prefix string) *Logger {
	return &Logger{
		loggingT: &logging,
		prefix:   prefix,
	}
}

// WithData creates a Logger with a given set of backend data.
// All args provided will be wrapped with glog.Data prior to being sent.
func WithData(args ...interface{}) *Logger {
	return &Logger{
		loggingT: &logging,
		data:     args,
	}
}

// WithPrefix creates a Logger from an existing logger with a specified prefix.
// Any prefix on the input Logger will be replaced.
func (l *Logger) WithPrefix(prefix string) *Logger {
	if l == nil {
		return nil
	}
	return &Logger{
		loggingT: l.loggingT,
		data:     l.data,
		prefix:   prefix,
	}
}

// WithData creates a Logger from an existing logger, using the specifed data.
// Any existing data on the input Logger will be replaced.
func (l *Logger) WithData(vars ...interface{}) *Logger {
	if l == nil {
		return nil
	}
	return &Logger{
		loggingT: l.loggingT,
		data:     vars,
		prefix:   l.prefix,
	}
}

// Info is equivalent to the global Info function, with the addition of prefix and data content from this Logger.
func (l *Logger) Info(args ...interface{}) {
	l.print(infoLog, l.extendWithPfx(args)...)
}

// InfoWithDepth is equivalent to the global InfoWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) InfoWithDepth(extraDepth int, args ...interface{}) {
	l.printWithDepth(infoLog, extraDepth, l.extendWithPfx(args)...)
}

// Infoln is equivalent to the global Infoln function, with the addition of prefix and data content from this Logger.
func (l *Logger) Infoln(args ...interface{}) {
	l.println(infoLog, l.extendWithPfx(args)...)
}

// InfolnWithDepth is equivalent to the global InfolnWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) InfolnWithDepth(extraDepth int, args ...interface{}) {
	l.printlnWithDepth(infoLog, extraDepth, l.extendWithPfx(args)...)
}

// Infof is equivalent to the global Infof function, with the addition of prefix and data content from this Logger.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.printf(infoLog, l.pfx(format), l.extend(args)...)
}

// InfofWithDepth is equivalent to the global InfofWithSepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) InfofWithDepth(extraDepth int, format string, args ...interface{}) {
	l.printfWithDepth(infoLog, extraDepth, l.pfx(format), l.extend(args)...)
}

// Warning is equivalent to the global Warning function, with the addition of prefix and data content from this Logger.
func (l *Logger) Warning(args ...interface{}) {
	l.print(warningLog, l.extendWithPfx(args)...)
}

// WarningWithDepth is equivalent to the global WarningWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) WarningWithDepth(extraDepth int, args ...interface{}) {
	l.printWithDepth(warningLog, extraDepth, l.extendWithPfx(args)...)
}

// Warningln is equivalent to the global Warningln function, with the addition of prefix and data content from this Logger.
func (l *Logger) Warningln(args ...interface{}) {
	l.println(warningLog, l.extendWithPfx(args)...)
}

// WarninglnWithDepth is equivalent to the global WarninglnWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) WarninglnWithDepth(extraDepth int, args ...interface{}) {
	l.printlnWithDepth(warningLog, extraDepth, l.extendWithPfx(args)...)
}

// Warningf is equivalent to the global Warningf function, with the addition of prefix and data content from this Logger.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.printf(warningLog, l.pfx(format), l.extend(args)...)
}

// WarningfWithDepth is equivalent to the global WarningfWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) WarningfWithDepth(extraDepth int, format string, args ...interface{}) {
	l.printfWithDepth(warningLog, extraDepth, l.pfx(format), l.extend(args)...)
}

// Error is equivalent to the global Error function, with the addition of prefix and data content from this Logger.
func (l *Logger) Error(args ...interface{}) {
	l.print(errorLog, l.extendWithPfx(args)...)
}

// ErrorWithDepth is equivalent to the global ErrorWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) ErrorWithDepth(extraDepth int, args ...interface{}) {
	l.printWithDepth(errorLog, extraDepth, l.extendWithPfx(args)...)
}

// Errorln is equivalent to the global Errorln function, with the addition of prefix and data content from this Logger.
func (l *Logger) Errorln(args ...interface{}) {
	l.println(errorLog, l.extendWithPfx(args)...)
}

// ErrorlnWithDepth is equivalent to the global ErrorlnWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) ErrorlnWithDepth(extraDepth int, args ...interface{}) {
	l.printlnWithDepth(errorLog, extraDepth, l.extendWithPfx(args)...)
}

// Errorf is equivalent to the global Errorf function, with the addition of prefix and data content from this Logger.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.printf(errorLog, l.pfx(format), l.extend(args)...)
}

// ErrorfWithDepth is equivalent to the global ErrorfWithDepth function, with the addition of prefix and data content from this Logger.
func (l *Logger) ErrorfWithDepth(extraDepth int, format string, args ...interface{}) {
	l.printfWithDepth(errorLog, extraDepth, l.pfx(format), l.extend(args)...)
}

func (l *Logger) extendWithPfx(args []interface{}) []interface{} {
	if l.prefix != "" {
		args = append([]interface{}{
			l.prefix,
		}, args...)
	}
	return l.extend(args)
}

func (l *Logger) extend(args []interface{}) []interface{} {
	for _, d := range l.data {
		args = append(args, data{d})
	}
	return args
}

func (l *Logger) pfx(log string) string {
	return fmt.Sprintf("%v %v", l.prefix, log)
}
