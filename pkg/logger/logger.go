package logger

import (
    "fmt"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "io"
    "os"
    "strings"
)

type LogLevel string

const (
    ErrorLevel LogLevel = "error"
    WarnLevel  LogLevel = "warn"
    InfoLevel  LogLevel = "info"
    DebugLevel LogLevel = "debug"
    PanicLevel LogLevel = "panic"
    FatalLevel LogLevel = "fatal"
)

// Interface -.
type Interface interface {
    Debug(message interface{}, args ...interface{})
    Info(message interface{}, args ...interface{})
    Warn(message interface{}, args ...interface{})
    Error(message interface{}, args ...interface{})
    Panic(message interface{}, args ...interface{})
    Fatal(message interface{}, args ...interface{})
    With(key Field, value interface{}) Interface
}

// Logger -.
type Logger struct {
    logger *zap.SugaredLogger
}

var _ Interface = (*Logger)(nil)

// New -.
func New(appName string, cnf Params, out io.Writer) *Logger {
    core := newZapCore(cnf, out)
    zap.NewProductionConfig()
    logger := zap.New(core)

    sugLogger := logger.Sugar()

    return &Logger{
        logger: sugLogger.With(AppName, appName),
    }
}

func toZapLevel(level LogLevel) zapcore.Level {
    switch LogLevel(strings.ToLower(string(level))) {
    case ErrorLevel:
        return zap.ErrorLevel
    case WarnLevel:
        return zap.WarnLevel
    case InfoLevel:
        return zap.InfoLevel
    case DebugLevel:
        return zap.DebugLevel
    case PanicLevel:
        return zap.PanicLevel
    case FatalLevel:
        return zap.FatalLevel
    default:
        return zap.InfoLevel
    }
}

type Params struct {
    AppName                  string
    LogDir                   string
    Level                    LogLevel
    UseStdAndFIle            bool
    AddLowPriorityLevelToCmd bool
}

func newZapCore(cnf Params, out io.Writer) (core zapcore.Core) {
    // First, define our level-handling logic.
    highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
        return lvl >= toZapLevel(cnf.Level)
    })

    if cnf.AddLowPriorityLevelToCmd { // separate levels
        lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
            return lvl < toZapLevel(cnf.Level)
        })

        topicDebugging := zapcore.AddSync(out)
        topicErrors := zapcore.AddSync(out)
        fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

        if cnf.UseStdAndFIle && cnf.LogDir != "" {
            consoleDebugging := zapcore.Lock(os.Stdout)
            consoleErrors := zapcore.Lock(os.Stderr)
            consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

            core = zapcore.NewTee(
                zapcore.NewCore(fileEncoder, topicErrors, highPriority),
                zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
                zapcore.NewCore(fileEncoder, topicDebugging, lowPriority),
                zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
            )
        } else {
            core = zapcore.NewTee(
                zapcore.NewCore(fileEncoder, topicErrors, highPriority),
                zapcore.NewCore(fileEncoder, topicDebugging, lowPriority),
            )
        }
    } else { // not separate levels
        topicErrors := zapcore.AddSync(out)
        fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())

        if cnf.UseStdAndFIle && cnf.LogDir != "" {
            consoleErrors := zapcore.Lock(os.Stderr)
            consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

            core = zapcore.NewTee(
                zapcore.NewCore(fileEncoder, topicErrors, highPriority),
                zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
            )
        } else {
            core = zapcore.NewTee(
                zapcore.NewCore(fileEncoder, topicErrors, highPriority),
            )
        }
    }
    return core
}

func (l *Logger) Sync() error {
    return l.logger.Sync()
}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
    l.log(l.logger.Debugf, message, args)
}

// Info -.
func (l *Logger) Info(message interface{}, args ...interface{}) {
    l.log(l.logger.Infof, message, args)
}

// Warn -.
func (l *Logger) Warn(message interface{}, args ...interface{}) {
    l.log(l.logger.Warnf, message, args)
}

// Panic -.
func (l *Logger) Panic(message interface{}, args ...interface{}) {
    l.log(l.logger.Panicf, message, args)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
    l.log(l.logger.Errorf, message, args)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
    l.log(l.logger.Fatalf, message, args)
}

func (l *Logger) log(lg func(message string, args ...interface{}), message interface{}, args ...interface{}) {
    switch tp := message.(type) {
    case error:
        lg(tp.Error(), args...)
    case string:
        lg(tp, args...)
    default:
        lg(fmt.Sprintf("message %v has unknown type %v", message, tp), args...)
    }
}

func (l *Logger) With(key Field, value interface{}) Interface {
    l.logger = l.logger.With(key, value)
    return l
}
