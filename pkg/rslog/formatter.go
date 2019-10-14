package rslog

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type Formatter struct{}

func (f Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	colorLevel := getColorLevel(entry.Level.String())
	logPrefix := fmt.Sprintf("%s\t%s\t", color.HiWhiteString(entry.Time.Format("2006-01-02 15:04:05.999")), colorLevel)
	fn, ok := entry.Data["func"]
	if ok {
		logPrefix += color.HiBlueString(fmt.Sprintf("%v\t", fn))
	}

	return []byte(fmt.Sprintf("%s: %s\n", logPrefix, color.HiWhiteString(entry.Message))), nil
}

func getColorLevel(lvl string) string {
	upperLevel := strings.ToUpper(lvl)
	switch upperLevel {
	case "PANIC", "FATAL", "ERROR":
		return color.RedString(upperLevel)
	case "WARN", "WARNING":
		return color.YellowString(upperLevel)
	case "INFO":
		return color.HiGreenString(upperLevel)
	case "DEBUG":
		return color.HiCyanString(upperLevel)
	case "TRACE":
		return color.MagentaString(upperLevel)
	default:
		return color.WhiteString(upperLevel)
	}
}
