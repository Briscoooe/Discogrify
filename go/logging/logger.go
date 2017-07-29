package logging

type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}