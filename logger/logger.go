package logger

import (
	"fmt"
	"log"
	"os"
)

const filename = "app.log"

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
)

func Init() error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("open log file failed: %w", err)
	}
	infoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return nil
}

func Info(msg string) {
	if infoLogger != nil {
		infoLogger.Output(2, msg)
	}
}

func Error(msg string) {
	if errorLogger != nil {
		errorLogger.Output(2, msg)

	}
}

func CloseFile() {
	if file != nil {
		file.Close()
	}

}
