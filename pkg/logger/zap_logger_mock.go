package logger

type mockLogger struct{}

func (m *mockLogger) InitLogger()                                  {}
func (m *mockLogger) Debug(args ...interface{})                    {}
func (m *mockLogger) Debugf(template string, args ...interface{})  {}
func (m *mockLogger) Info(args ...interface{})                     {}
func (m *mockLogger) Infof(template string, args ...interface{})   {}
func (m *mockLogger) Warn(args ...interface{})                     {}
func (m *mockLogger) Warnf(template string, args ...interface{})   {}
func (m *mockLogger) Error(args ...interface{})                    {}
func (m *mockLogger) Errorf(template string, args ...interface{})  {}
func (m *mockLogger) DPanic(args ...interface{})                   {}
func (m *mockLogger) DPanicf(template string, args ...interface{}) {}
func (m *mockLogger) Fatal(args ...interface{})                    {}
func (m *mockLogger) Fatalf(template string, args ...interface{})  {}

func NewMockLogger() Logger {
	return &mockLogger{}
}
