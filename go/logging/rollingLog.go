package logging

import (
	"log"
	"os"
	"fmt"
	"github.com/natefinch/lumberjack"
)

type RollingLogger struct {
	RollingLog *log.Logger
}

func NewRollingLogger(filename string, maxSize, maxBackups, maxAge int) *RollingLogger {
	e, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(1)
	}
	rollingLog := log.New(e, "", log.Ldate|log.Ltime)
	rollingLog.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,   // megabytes after which new file is created
		MaxBackups: maxBackups, // number of backups
		MaxAge:     maxAge,  //days
	})
	return &RollingLogger{
		RollingLog: rollingLog,
	}
}

func (r RollingLogger) Fatal(v ...interface{}) {
	fmt.Println(v...)
	r.RollingLog.Println(v...)
}

func (r RollingLogger) Fatalf(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v...))
	r.RollingLog.Printf(fmt.Sprintf(format, v...))
}

func (r RollingLogger) Printf(format string, v ...interface{}) {
	fmt.Println(fmt.Sprintf(format, v...))
	r.RollingLog.Printf(fmt.Sprintf(format, v...))
}

func (r RollingLogger) Println(v ...interface{}) {
	fmt.Println(v...)
	r.RollingLog.Println(v...)
}
