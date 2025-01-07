package testutils

import (
	"fmt"
)

type MockLogger struct {
	Logs []string
}

func (m *MockLogger) Printf(format string, v ...interface{}) {
	m.Logs = append(m.Logs, fmt.Sprintf(format, v...))
}
