package logger

import (
	"log"
	"log/syslog"
)

type Level int

const (
	_ Level = iota
	DebugLvl
	InfoLvl
	WarningLvl
	ErrorLvl
	FatalLvl
)

type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	fatal   *log.Logger
}

func NewLogger(level Level) (*Logger, error) {
	// Создаем логгеры для разных уровней
	debug, err := syslog.NewLogger(syslog.LOG_DEBUG, log.LstdFlags)
	if err != nil {
		return nil, err
	}

	info, err := syslog.NewLogger(syslog.LOG_INFO, log.LstdFlags)
	if err != nil {
		return nil, err
	}

	warning, err := syslog.NewLogger(syslog.LOG_WARNING, log.LstdFlags)
	if err != nil {
		return nil, err
	}

	errLog, err := syslog.NewLogger(syslog.LOG_ERR, log.LstdFlags)
	if err != nil {
		return nil, err
	}

	fatal, err := syslog.NewLogger(syslog.LOG_CRIT, log.LstdFlags)
	if err != nil {
		return nil, err
	}

	// Устанавливаем уровень детализации
	logger := &Logger{
		debug:   debug,
		info:    info,
		warning: warning,
		error:   errLog,
		fatal:   fatal,
	}

	// Включаем/отключаем уровни в зависимости от настройки
	switch level {
	case DebugLvl:
		logger.debug = debug
	case InfoLvl:
		logger.debug = nil
	case WarningLvl:
		logger.debug = nil
		logger.info = nil
	case ErrorLvl:
		logger.debug = nil
		logger.info = nil
		logger.warning = nil
	case FatalLvl:
		logger.debug = nil
		logger.info = nil
		logger.warning = nil
		logger.error = nil
	default:
		logger.debug = nil
	}

	return logger, nil
}

// Debug логирует сообщение с уровнем DEBUG
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.debug != nil {
		l.debug.Printf(format, v...)
	}
}

// Info логирует сообщение с уровнем INFO
func (l *Logger) Info(format string, v ...interface{}) {
	if l.info != nil {
		l.info.Printf(format, v...)
	}
}

// Warning логирует сообщение с уровнем WARNING
func (l *Logger) Warning(format string, v ...interface{}) {
	if l.warning != nil {
		l.warning.Printf(format, v...)
	}
}

// Error логирует сообщение с уровнем ERROR
func (l *Logger) Error(format string, v ...interface{}) {
	if l.error != nil {
		l.error.Printf(format, v...)
	}
}

// Fatal логирует сообщение с уровнем FATAL и завершает программу
func (l *Logger) Fatal(format string, v ...interface{}) {
	if l.fatal != nil {
		l.fatal.Fatalf(format, v...)
	}
}
