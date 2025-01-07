package interfaces

type Logger interface {
	Printf(format string, v ...interface{})
}
