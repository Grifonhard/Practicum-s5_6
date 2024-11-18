package logger

import (
	"fmt"
)

var log *Logger

func InitLogger(level string, logFilePath string) error {
	lognew, err := new(level, logFilePath)
	if err != nil {
		return err
	}
	log = lognew
	return nil
}

func Info(msg string, a ...any) {
	log.info(fmt.Sprintf(msg, a...))
}

func Debug(msg string, a ...any) {
	log.debug(fmt.Sprintf(msg, a...))
}

func Warn(msg string, a ...any) {
	log.warn(fmt.Sprintf(msg, a...))
}

func Error(msg string, a ...any) {
	log.error(fmt.Sprintf(msg, a...))
}

func Fatal(msg string, a ...any) {
	log.fatal(fmt.Sprintf(msg, a...))
}

func Sync() {
	log.sync()
}
