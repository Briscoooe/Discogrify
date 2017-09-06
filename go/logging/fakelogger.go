package logging

type FakeLogger struct {
}

func (f FakeLogger) Fatal(v ...interface{}) {
}
func (f FakeLogger) Fatalf(format string, v ...interface{}) {
}
func (f FakeLogger) Println(v ...interface{}) {
}
func (f FakeLogger) Printf(format string, v ...interface{}) {
}

func NewFakeLogger() FakeLogger {
	return FakeLogger{}
}