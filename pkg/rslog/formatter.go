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
	logPrefix := fmt.Sprintf("[%s][%s]", colorLevel, color.HiWhiteString(entry.Time.Format("2006-01-02T15:04:05.999999-07:00")))
	for k, v := range entry.Data {
		logPrefix += color.BlueString(fmt.Sprintf("[%s:%v]", k, v))
	}

	return []byte(fmt.Sprintf("%s :: %s\n", logPrefix, entry.Message)), nil
}

func getColorLevel(lvl string) string {
	upperLevel := strings.ToUpper(lvl)
	switch upperLevel {
	case "PANIC", "FATAL", "ERROR":
		return color.RedString(upperLevel)
	case "WARN", "WARNING":
		return color.YellowString(upperLevel)
	case "INFO":
		return color.GreenString(upperLevel)
	case "DEBUG":
		return color.CyanString(upperLevel)
	case "TRACE":
		return color.MagentaString(upperLevel)
	default:
		return color.WhiteString(upperLevel)
	}
}
