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
	return &Logger{
		loggingT: l.loggingT,
		data:     l.data,
		prefix:   prefix,
	}
}

// WithData creates a Logger from an existing logger, using the specifed data.
// Any existing data on the input Logger will be replaced.
func (l *Logger) WithData(vars ...interface{}) *Logger {
	return &Logger{
		loggingT: l.loggingT,
		data:     vars,
		prefix:   l.prefix,
	}
}

// AppendData creates a Logger from an existing logger,
// appending the provided data to the data in the existing logger.
func (l *Logger) AppendData(vars ...interface{}) *Logger {
	newData := make([]interface{}, len(l.data))
	copy(newData, l.data)
	return &Logger{
		loggingT: l.loggingT,
		data:     append(newData, vars),
		prefix:   l.prefix,
	}
}

// Info is equivalent to the global Info function, with the addition of prefix and data content from this Logger.
func (l *Logger) Info(args ...interface{}) {
	l.print(infoLog, l.extendWithPfx(args)...)
}

// Infoln is equivalent to the global Infoln function, with the addition of prefix and data content from this Logger.
func (l *Logger) Infoln(args ...interface{}) {
	l.println(infoLog, l.extendWithPfx(args)...)
}

// Infof is equivalent to the global Infof function, with the addition of prefix and data content from this Logger.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.printf(infoLog, l.pfx(format), l.extend(args)...)
}

// Warning is equivalent to the global Warning function, with the addition of prefix and data content from this Logger.
func (l *Logger) Warning(args ...interface{}) {
	l.print(warningLog, l.extendWithPfx(args)...)
}

// Warningln is equivalent to the global Warningln function, with the addition of prefix and data content from this Logger.
func (l *Logger) Warningln(args ...interface{}) {
	l.println(warningLog, l.extendWithPfx(args)...)
}

// Warningf is equivalent to the global Warningf function, with the addition of prefix and data content from this Logger.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.printf(warningLog, l.pfx(format), l.extend(args)...)
}

// Error is equivalent to the global Error function, with the addition of prefix and data content from this Logger.
func (l *Logger) Error(args ...interface{}) {
	l.print(errorLog, l.extendWithPfx(args)...)
}

// GetErrorEvent is equivalent to the global GetErrorEvent function, with the addition of prefix and data content from this Logger.
func (l *Logger) GetErrorEvent(args ...interface{}) Event {
	return l.getEvent(errorLog, l.extendWithPfx(args))
}

// ErrorIf is equivalent to the global ErrorIf function, with the addition of prefix and data content from this Logger.
func (l *Logger) ErrorIf(err error, args ...interface{}) {
	if err != nil {
		if args != nil {
			args = append(args, ": ")
		}
		args = append(args, err)
		l.print(errorLog, l.extendWithPfx(args)...)
	}
}

// Errorln is equivalent to the global Errorln function, with the addition of prefix and data content from this Logger.
func (l *Logger) Errorln(args ...interface{}) {
	l.println(errorLog, l.extendWithPfx(args)...)
}

// Errorf is equivalent to the global Errorf function, with the addition of prefix and data content from this Logger.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.printf(errorLog, l.pfx(format), l.extend(args)...)
}

// ErrorfIf is equivalent to the global ErrorfIf function, with the addition of prefix and data content from this Logger.
func (l *Logger) ErrorfIf(err error, format string, args ...interface{}) {
	if err != nil {
		format += ": %v"
		args = append(args, err)
		l.printf(errorLog, l.pfx(format), l.extend(args)...)
	}
}

// Fatal is equivalent to the global Fatal  function, with the addition of prefix and data content from this Logger.
func (l *Logger) Fatal(args ...interface{}) {
	l.print(fatalLog, l.extendWithPfx(args)...)
}

// FatalIf is equivalent to the global FatalIf  function, with the addition of prefix and data content from this Logger.
func (l *Logger) FatalIf(err error, args ...interface{}) {
	if err != nil {
		if args != nil {
			errStr := ": " + err.Error()
			args = append(args, errStr)
			l.print(fatalLog, l.extendWithPfx(args)...)
		} else {
			l.print(fatalLog, l.extendWithPfx([]interface{}{err}))
		}
	}
}

// Fatalln is equivalent to the global Fatalln  function, with the addition of prefix and data content from this Logger.
func (l *Logger) Fatalln(args ...interface{}) {
	l.println(fatalLog, l.extendWithPfx(args)...)
}

// Fatalf is equivalent to the global Fatalf  function, with the addition of prefix and data content from this Logger.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.printf(fatalLog, l.pfx(format), l.extend(args)...)
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
		args = append(args, Data(d))
	}
	return args
}

func (l *Logger) pfx(log string) string {
	return fmt.Sprintf("%v %v", l.prefix, log)
}
