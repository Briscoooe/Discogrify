package logging

type FakeLogger struct {
}

func (r FakeLogger) LogErr(err error, v ...interface{}) {
}
func (r FakeLogger) LogErrf(err error, format string, v ... interface{}) {
}
func (f FakeLogger) Log(v ...interface{}) {
}
func (f FakeLogger) Logf(format string, v ...interface{}) {
}
func (f FakeLogger) Fatal(v ...interface{}) {
}
func (f FakeLogger) Fatalf(format string, v ...interface{}) {
}

func NewFakeLogger() FakeLogger {
	return FakeLogger{}
}