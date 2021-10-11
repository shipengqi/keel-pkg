package log

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/shipengqi/keel-pkg/lib/spinner"
	"github.com/shipengqi/keel-pkg/lib/utils/fmtutil"
	"github.com/shipengqi/keel-pkg/lib/utils/fsutil"
)

// Config Configuration for logging
type Config struct {
	ConsoleEnabled bool
	ConsoleJson    bool

	FileJson    bool
	FileEnabled bool

	ConsoleLevel string
	FileLevel    string

	// Directory to log when FileEnabled is true
	Directory string
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string
}

type Interface interface {
	Debugt(msg string, fields ...zapcore.Field)
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Debugs(args ...interface{})

	Infot(msg string, fields ...zapcore.Field)
	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Infos(args ...interface{})

	Warnt(msg string, fields ...zapcore.Field)
	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Warns(args ...interface{})

	Errort(msg string, fields ...zapcore.Field)
	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Errors(args ...interface{})

	Panict(msg string, fields ...zapcore.Field)
	Panicf(template string, args ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Panic(msg string, keysAndValues ...interface{})
	Panics(args ...interface{})

	Fatalt(msg string, fields ...zapcore.Field)
	Fatalf(template string, args ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	Fatals(args ...interface{})

	AtLevel(level zapcore.Level, msg string, fields ...zapcore.Field) *Logger
}

type Logger struct {
	Unsugared *zap.Logger
	*zap.SugaredLogger
}

func (l *Logger) AtLevel(level zapcore.Level, msg string, fields ...zapcore.Field) *Logger {
	switch level {
	case zapcore.DebugLevel:
		l.Unsugared.Debug(msg, fields...)
	case zapcore.PanicLevel:
		l.Unsugared.Panic(msg, fields...)
	case zapcore.ErrorLevel:
		l.Unsugared.Error(msg, fields...)
	case zapcore.WarnLevel:
		l.Unsugared.Warn(msg, fields...)
	case zapcore.InfoLevel:
		l.Unsugared.Info(msg, fields...)
	case zapcore.FatalLevel:
		l.Unsugared.Fatal(msg, fields...)
	default:
		l.Unsugared.Warn("Logging at unknown level", zap.Any("level", level))
		l.Unsugared.Warn(msg, fields...)
	}
	return l
}

func (l *Logger) Debugt(msg string, fields ...zapcore.Field) {
	l.Unsugared.Debug(msg, fields...)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Debugs(args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}

func (l *Logger) Infot(msg string, fields ...zapcore.Field) {
	l.Unsugared.Info(msg, fields...)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

func (l *Logger) Infos(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func (l *Logger) Warnt(msg string, fields ...zapcore.Field) {
	l.Unsugared.Warn(msg, fields...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Warns(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

func (l *Logger) Errort(msg string, fields ...zapcore.Field) {
	l.Unsugared.Error(msg, fields...)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

func (l *Logger) Errors(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

func (l *Logger) Panict(msg string, fields ...zapcore.Field) {
	l.Unsugared.Panic(msg, fields...)
}

func (l *Logger) Panic(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Panicw(msg, keysAndValues...)
}

func (l *Logger) Panics(args ...interface{}) {
	l.SugaredLogger.Panic(args...)
}

func (l *Logger) Fatalt(msg string, fields ...zapcore.Field) {
	l.Unsugared.Fatal(msg, fields...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Fatals(args ...interface{}) {
	l.SugaredLogger.Fatal(args...)
}

// How to log, by example:
// log.Infot("Importing new file", zap.String("source", filename), zap.Int("size", 1024))
// log.Info("Importing new file", "source", filename, "size", 1024)
// To log a stacktrace:
// log.Errort("It went wrong, zap.Stack())

// defaultZapLogger is the default logger instance that should be used to log
// It's assigned a default value here for tests (which do not call log.Configure())
var defaultZapLogger *Logger

func init() {
	defaultZapLogger, _ = newZapLogger(Config{
		ConsoleEnabled: true,
		ConsoleLevel:   "INFO",
	})
}

func Debugt(msg string, fields ...zapcore.Field) {
	defaultZapLogger.Debugt(msg, fields...)
}

func Debugf(template string, args ...interface{}) {
	defaultZapLogger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Debugw(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Debug(msg, keysAndValues...)
}

func Debugs(args ...interface{}) {
	defaultZapLogger.Debugs(args...)
}

func Infot(msg string, fields ...zapcore.Field) {
	defaultZapLogger.Infot(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	defaultZapLogger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Infow(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Info(msg, keysAndValues...)
}

func Infos(args ...interface{}) {
	defaultZapLogger.Infos(args...)
}

func Warnt(msg string, fields ...zapcore.Field) {
	defaultZapLogger.Warnt(msg, fields...)
}

func Warnf(template string, args ...interface{}) {
	defaultZapLogger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Warnw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Warn(msg, keysAndValues...)
}

func Warns(args ...interface{}) {
	defaultZapLogger.Warns(args...)
}

func Errort(msg string, fields ...zapcore.Field) {
	defaultZapLogger.Errort(msg, fields...)
}

func Errorf(template string, args ...interface{}) {
	defaultZapLogger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Errorw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Error(msg, keysAndValues...)
}

func Errors(args ...interface{}) {
	defaultZapLogger.Errors(args...)
}

func Panict(msg string, fields ...zapcore.Field) {
	defaultZapLogger.Panict(msg, fields...)
}

func Panicf(template string, args ...interface{}) {
	defaultZapLogger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Panicw(msg, keysAndValues...)
}

func Panic(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Panic(msg, keysAndValues...)
}

func Panics(args ...interface{}) {
	defaultZapLogger.Panics(args...)
}

func Fatalt(msg string, fields ...zapcore.Field) {
	defaultZapLogger.Fatalt(msg, fields...)
}

func Fatalf(template string, args ...interface{}) {
	defaultZapLogger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Fatalw(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...interface{}) {
	defaultZapLogger.Fatal(msg, keysAndValues...)
}

func Fatals(args ...interface{}) {
	defaultZapLogger.Fatals(args...)
}

// AtLevel logs the message at a specific log level
func AtLevel(level zapcore.Level, msg string, fields ...zapcore.Field) {
	switch level {
	case zapcore.DebugLevel:
		Debugt(msg, fields...)
	case zapcore.PanicLevel:
		Panict(msg, fields...)
	case zapcore.ErrorLevel:
		Errort(msg, fields...)
	case zapcore.WarnLevel:
		Warnt(msg, fields...)
	case zapcore.InfoLevel:
		Infot(msg, fields...)
	case zapcore.FatalLevel:
		Fatalt(msg, fields...)
	default:
		Warnt("Logging at unknown level", zap.Any("level", level))
		Warnt(msg, fields...)
	}
}

func Progress(msg, status string) {
	pretty := fmtutil.Pretty(msg, ColorizeStatus(status))
	Info(pretty)
}

func ProgressSub(msg, status string) {
	pretty := fmtutil.Pretty(msg, ColorizeStatus(status))
	Info(pretty)
}

var defaultSpinner = spinner.New()
// Sprogress start a progress, must call the StopProgress after the progress end
func Sprogress(msg string) *spinner.Spinner {
	prefix := fmtutil.PrettyPrefix(msg)
	defaultSpinner.Reset().WithPrefix(prefix).WithSuffix(" ]").Start()
	return defaultSpinner
}

// StopProgress stop a progress
func StopProgress(status string) {
	endStr := defaultSpinner.StopWithStatus(ColorizeStatus(status))
	Debug(endStr)
}

// Configure sets up the logging framework
func Configure(config Config) (*Logger, error) {
	logger, err := newZapLogger(config)
	if err != nil {
		return nil, err
	}
	logger.Debugt("logging configured",
		zap.Bool("consoleEnabled", config.ConsoleEnabled),
		zap.String("consoleLevel", config.ConsoleLevel),
		zap.Bool("consoleJson", config.ConsoleJson),
		zap.Bool("fileEnabled", config.FileEnabled),
		zap.String("fileLevel", config.FileLevel),
		zap.Bool("fileJson", config.FileJson),
		zap.String("logDirectory", config.Directory),
		zap.String("fileName", config.Filename))

	defaultZapLogger = logger
	zap.RedirectStdLog(defaultZapLogger.Unsugared)
	return logger, nil
}

func newZapLogger(config Config) (*Logger, error) {
	if !config.ConsoleEnabled && !config.FileEnabled {
		return nil, errors.New("no logger enabled")
	}

	var cores []zapcore.Core
	var consoleLevel zapcore.Level
	var fileLevel zapcore.Level
	var err error
	if config.ConsoleEnabled {
		err = consoleLevel.Set(strings.ToLower(config.ConsoleLevel))
		if err != nil {
			return nil, errors.Wrap(err, "set console level")
		}
		consoleEncCfg := zapcore.EncoderConfig{
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
		}
		consoleLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= consoleLevel
		})
		consoleEncoder := zapcore.NewConsoleEncoder(consoleEncCfg)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), consoleLevelEnabler))
	}

	if config.FileEnabled {
		err = fileLevel.Set(strings.ToLower(config.FileLevel))
		if err != nil {
			return nil, errors.Wrap(err, "set file level")
		}

		jsonEncCfg := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.NanosDurationEncoder,
		}
		fileLevelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= fileLevel
		})
		fileEncoder := zapcore.NewJSONEncoder(jsonEncCfg)
		err = fsutil.MustMkDir(config.Directory)
		if err != nil {
			return nil, errors.Wrap(err, "MustMkDir")
		}
		file := filepath.Join(config.Directory, config.Filename)
		fd, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
		if err != nil {
			return nil, errors.Wrap(err, "OpenFile")
		}
		cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(fd), fileLevelEnabler))
	}
	core := zapcore.NewTee(cores...)
	unsugared := zap.New(core)
	return &Logger{
		Unsugared: unsugared,
		SugaredLogger: unsugared.Sugar(),
	}, nil
}

// ColorizeStatus a status string based on given color.
func ColorizeStatus(s string) string {
	var color fmtutil.Color
	switch strings.ToLower(s) {
	case "ok", "pass", "success":
		color = fmtutil.FgGreen
		break
	case "failed", "error":
		color = fmtutil.FgRed
		break
	case "warn", "warning", "skip":
		color = fmtutil.FgYellow
		break
	default:
		color = fmtutil.FgWhite
		break
	}
	return fmtutil.Colorize(s, color)
}
