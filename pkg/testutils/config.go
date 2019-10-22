package testutils

import (
	"github.com/realsangil/apimonitor/pkg/rslog"
)

type logConfig string

func (l logConfig) GetLevel() string {
	return rslog.LevelDebug
}

func (l logConfig) GetFormat() string {
	return ""
}

func (l logConfig) GetOutput() string {
	return rslog.OutputConsole
}

func (l logConfig) GetPath() string {
	return ""
}

func SetLogConfig() {
	c := logConfig("")
	rslog.Init(&c)
}
