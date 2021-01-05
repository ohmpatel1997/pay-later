package log

import (
	"io"
	"os"
	"runtime"
	"runtime/debug"

	lr "github.com/sirupsen/logrus"
)

const (
	stackKey          = "stack"
	reportLocationKey = "reportLocation"
)

type Logger interface {
	Debug(string, ...interface{}) // Debug("")
	DebugD(string, Fields)

	Info(string, ...interface{})
	InfoD(string, Fields)

	Warn(string, ...interface{})
	WarnD(string, Fields)

	Error(string, ...interface{})
	ErrorD(string, Fields)

	Panic(string, ...interface{})
	PanicD(string, Fields)
}

type logger struct {
	l *lr.Logger
}

type Fields lr.Fields

func NewLogger(opts ...Option) Logger {

	la := &loggerOpts{
		output: os.Stdout,
		format: &lr.JSONFormatter{},
	}

	for _, opt := range opts {
		opt(la)
	}

	l := lr.New()

	l.SetOutput(la.output)
	l.SetFormatter(la.format)

	if os.Getenv("API_ENV") == "prod" {
		l.SetLevel(lr.InfoLevel)
	} else {
		l.SetLevel(lr.DebugLevel)
	}

	return logger{
		l: l,
	}
}

type loggerOpts struct {
	output io.Writer
	format lr.Formatter
}

type Option func(*loggerOpts)

func SetOutput(i io.Writer) Option {
	return func(opts *loggerOpts) {
		opts.output = i
	}
}

func SetFormat(f lr.Formatter) Option {
	return func(opts *loggerOpts) {
		opts.format = f
	}
}

func (l logger) Debug(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Debug(s)
}

func (l logger) DebugD(s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Debug(s)
}

func (l logger) Info(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Info(s)
}

func (l logger) InfoD(s string, f Fields) {
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Info(s)
}

func (l logger) Warn(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Warn(s)
}

func (l logger) WarnD(s string, f Fields) {
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Warn(s)
}

func (l logger) Error(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Error(s)
}

func (l logger) ErrorD(s string, f Fields) {
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Error(s)
}

func (l logger) Panic(s string, fs ...interface{}) {
	f := getFields(fs)
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Panic(s)
}

func (l logger) PanicD(s string, f Fields) {
	f = l.appendStack(f)
	f = l.appendReportLocation(f)
	l.l.WithFields(f.format()).Panic(s)
}

func (l logger) appendStack(f Fields) Fields {
	f[stackKey] = string(debug.Stack())

	return f
}

func (l logger) appendReportLocation(f Fields) Fields {
	pc, fn, line, _ := runtime.Caller(2)

	f[reportLocationKey] = map[string]interface{}{
		"filePath":     fn,
		"line":         line,
		"functionName": runtime.FuncForPC(pc).Name(),
	}

	return f
}

func (f Fields) format() lr.Fields {
	return lr.Fields(f)
}

func getFields(vfs ...interface{}) Fields {
	f := Fields{}

	if len(vfs) == 0 {
		return f
	}

	fs := vfs[0].([]interface{})

	if len(fs) > 0 {
		for i := 0; i < len(fs); i = i + 2 {
			if len(fs) <= i+1 {
				//If we have an odd number of args this case may pop up, just break
				break
			}

			key, ok := fs[i].(string)
			if !ok {
				break
			}

			value := fs[i+1]

			f[key] = value
		}
	}

	return f
}
