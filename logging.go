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
	l.File.Println(msg)
	l.Println(msg)
}
func (l *Logger) Record(str string, args ...interface{}) {
	msg := fmt.Sprintf(str, args...)
	l.File.Println(msg)
	l.Println(msg)
}

func (l *Logger) NewLine() {
	fmt.Println("")
}

func (l *Logger) Inform(args ...interface{}) {
	l.Logger.Println(args...)
}

type Logger struct {
	File *log.Logger
	*log.Logger
}

func NewLogger() (*Logger, error) {
	f, err := os.OpenFile(LOG_FILE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	l := new(Logger)
	l.File = log.New(f, "", log.Ldate|log.Ltime)
	l.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return l, nil
}
