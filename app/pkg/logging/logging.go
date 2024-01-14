package logging

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"time"
)

type logger struct {
	*logrus.Logger
}

type Logger interface {
	SetLevel(level logrus.Level)
	GetLevel() logrus.Level
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
	WithError(err error) *logrus.Entry
	WithContext(ctx context.Context) *logrus.Entry
	WithTime(t time.Time) *logrus.Entry
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

func GetLogger(ctx context.Context) Logger {
	return loggerFromContext(ctx)
}

func NewLogger() Logger {
	l := logrus.New()
	l.SetLevel(logrus.InfoLevel)
	// format logs settings
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s:%d", filename, f.Line), fmt.Sprintf("%s()", f.Function)
		},
		DisableColors: true,
		FullTimestamp: true,
	}

	l.SetOutput(os.Stdout)

	return &logger{
		Logger: l,
	}
}

func (l *logger) SetLevel(level logrus.Level) {
	l.Logger.SetLevel(level)
}

func (l *logger) GetLevel() logrus.Level {
	return l.Logger.GetLevel()
}

func (l *logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

func (l *logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

func (l *logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

func (l *logger) WithContext(ctx context.Context) *logrus.Entry {
	return l.Logger.WithContext(ctx)
}

func (l *logger) WithTime(t time.Time) *logrus.Entry {
	return l.Logger.WithTime(t)
}

func (l *logger) Trace(args ...interface{}) {
	l.Logger.Traceln(args...)
}

func (l *logger) Debug(args ...interface{}) {
	l.Logger.Debugln(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.Logger.Infoln(args...)
}

func (l *logger) Print(args ...interface{}) {
	l.Logger.Println(args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.Logger.Warningln(args...)
}

func (l *logger) Error(args ...interface{}) {
	l.Logger.Errorln(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.Logger.Fatalln(args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.Logger.Panicln(args...)
}
