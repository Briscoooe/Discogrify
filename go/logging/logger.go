package logging

type Logger interface {
	LogErr(err error, v ...interface{})
	LogErrf(err error, format string, v ... interface{})
	Log(v ...interface{})
	Logf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}