package rslog

type LogConfig interface {
	GetLevel() string
	GetFormat() string
	GetOutput() string
	GetPath() string
}
