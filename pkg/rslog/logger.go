package rslog

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/realsangil/apimonitor/pkg/rserrors"
)

const (
	FormatJSON = "json"
	FormatText = "text"

	OutputFile    = "file"
	OutputConsole = "console"

	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWanr  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

var (
	std = logrus.New()
)

func init() {
	std.SetFormatter(&Formatter{})
}

func setLevel(level string) error {
	var logLevel logrus.Level
	switch level {
	case LevelDebug:
		logLevel = logrus.DebugLevel
	case LevelInfo:
		logLevel = logrus.InfoLevel
	case LevelWanr:
		logLevel = logrus.WarnLevel
	case LevelError:
		logLevel = logrus.ErrorLevel
	case LevelFatal:
		logLevel = logrus.FatalLevel
	default:
		return rserrors.Error("invalid logger level")
	}
	std.SetLevel(logLevel)
	return nil
}

func setOutput(output, path string) error {
	switch output {
	case OutputFile:
		lumberjackLogger := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    100,
			MaxBackups: 30,
			MaxAge:     30,
			LocalTime:  true,
			Compress:   true,
		}
		std.SetOutput(lumberjackLogger)
	case OutputConsole:
	default:
		return rserrors.Error("invalid logger output")
	}
	return nil
}

func Init(config LogConfig) error {
	if err := setLevel(config.GetLevel()); err != nil {
		return errors.WithStack(err)
	}

	if err := setOutput(config.GetOutput(), config.GetPath()); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Trace(args ...interface{}) {
	getStdWithPrefix().Trace(args...)
}

func Debug(args ...interface{}) {
	getStdWithPrefix().Debug(args...)
}

func Print(args ...interface{}) {
	getStdWithPrefix().Print(args...)
}

func Info(args ...interface{}) {
	getStdWithPrefix().Info(args...)
}

func Warn(args ...interface{}) {
	getStdWithPrefix().Warn(args...)
}

func Warning(args ...interface{}) {
	getStdWithPrefix().Warning(args...)
}

func Error(args ...interface{}) {
	getStdWithPrefix().Error(args...)
}

func Panic(args ...interface{}) {
	getStdWithPrefix().Panic(args...)
}

func Fatal(args ...interface{}) {
	getStdWithPrefix().Fatal(args...)
}

func Tracef(format string, args ...interface{}) {
	getStdWithPrefix().Tracef(format, args...)
}

func Debugf(format string, args ...interface{}) {
	getStdWithPrefix().Debugf(format, args...)
}

func Printf(format string, args ...interface{}) {
	getStdWithPrefix().Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	getStdWithPrefix().Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	getStdWithPrefix().Warnf(format, args...)
}

func Warningf(format string, args ...interface{}) {
	getStdWithPrefix().Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	getStdWithPrefix().Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	getStdWithPrefix().Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	getStdWithPrefix().Fatalf(format, args...)
}

func Traceln(args ...interface{}) {
	getStdWithPrefix().Traceln(args...)
}

func Debugln(args ...interface{}) {
	getStdWithPrefix().Debugln(args...)
}

func Println(args ...interface{}) {
	getStdWithPrefix().Println(args...)
}

func Infoln(args ...interface{}) {
	getStdWithPrefix().Infoln(args...)
}

func Warnln(args ...interface{}) {
	getStdWithPrefix().Warnln(args...)
}

func Warningln(args ...interface{}) {
	getStdWithPrefix().Warningln(args...)
}

func Errorln(args ...interface{}) {
	getStdWithPrefix().Errorln(args...)
}

func Panicln(args ...interface{}) {
	getStdWithPrefix().Panicln(args...)
}

func Fatalln(args ...interface{}) {
	getStdWithPrefix().Fatalln(args...)
}

func getStdWithPrefix() *logrus.Entry {
	pc := make([]uintptr, 15)
	n := runtime.Callers(3, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return std.WithFields(logrus.Fields{
		"func": fmt.Sprintf("%s:%d", frame.Function, frame.Line),
	})
}
