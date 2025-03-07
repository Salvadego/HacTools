package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelError
	LevelNone
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger

	currentLevel = LevelError
)

func init() {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	InfoLogger = log.New(os.Stdout, "INFO: ", flags)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", flags)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", flags)

	SetLogLevel(currentLevel)
}

func SetLogLevel(level LogLevel) {
	currentLevel = level

	switch level {
	case LevelNone:
		InfoLogger.SetOutput(io.Discard)
		DebugLogger.SetOutput(io.Discard)
		ErrorLogger.SetOutput(io.Discard)
	case LevelError:
		InfoLogger.SetOutput(io.Discard)
		DebugLogger.SetOutput(io.Discard)
		ErrorLogger.SetOutput(os.Stderr)
	case LevelInfo:
		InfoLogger.SetOutput(os.Stdout)
		DebugLogger.SetOutput(io.Discard)
		ErrorLogger.SetOutput(os.Stderr)
	case LevelDebug:
		InfoLogger.SetOutput(os.Stdout)
		DebugLogger.SetOutput(os.Stdout)
		ErrorLogger.SetOutput(os.Stderr)
	}
}

func LogLevelFromString(level string) LogLevel {
	switch level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "error":
		return LevelError
	case "none":
		return LevelNone
	default:
		return LevelError
	}
}

func Info(format string, v ...any) {
	if currentLevel <= LevelInfo {
		InfoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...any) {
	if currentLevel <= LevelError {
		ErrorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debug(format string, v ...any) {
	if currentLevel <= LevelDebug {
		DebugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}
