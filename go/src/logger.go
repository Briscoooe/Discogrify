package main

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"log"
	"net/http"
	"os"
	"time"
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
func logRouter(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		rollingLog.Printf(
			r.Method,
			r.RequestURI,
			r.Host,
			name,
			time.Since(start),
		)
	})
}
