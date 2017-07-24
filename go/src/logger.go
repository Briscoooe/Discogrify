package main

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"log"
	"os"
)

var (
	rollingLog *log.Logger
)

func setupLogger() {
	e, err := os.OpenFile("./discogrify.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Printf("Error opening file: %v", err)
		os.Exit(1)
	}
	rollingLog = log.New(e, "", log.Ldate|log.Ltime)
	rollingLog.SetOutput(&lumberjack.Logger{
		Filename:   "./discogrify.log",
		MaxSize:    1,   // megabytes after which new file is created
		MaxBackups: 100, // number of backups
		MaxAge:     28,  //days
	})
}