package main

import (
	"fmt"
	"log"
	"os"
)

const LOG_FILE_NAME = FILE_DIR_NAME + "LOG_FILE.txt"

// Stdout logging may be replaced with file-logging, so
// just creating a simple wrapper func for now
func (l *Logger) ServerErr(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	l.Err.Println(msg)
	l.Println(msg)
}

func (l *Logger) NewLine() {
	fmt.Println("")
}

func (l *Logger) Inform(args ...interface{}) {
	l.Logger.Println(args...)
}

type Logger struct {
	Err *log.Logger
	*log.Logger
}

func NewLogger() (*Logger, error) {
	f, err := os.OpenFile(LOG_FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	l := new(Logger)
	l.Err = log.New(f, "", log.Lshortfile|log.Ldate|log.Ltime)
	l.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return l, nil
}
