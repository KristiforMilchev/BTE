package interfaces

type ILogger interface {
	Write(p []byte) (n int, err error)
	Logs() []string
	Close() error
}
