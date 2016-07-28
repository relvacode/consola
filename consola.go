// consola is a simple package for printing better looking Logrus log messages to the console.
package consola

import (
	"bytes"
	"fmt"
	"github.com/Sirupsen/logrus"
	"strings"
)

// Terminal color constants
const (
	Red      = "\x1b[31m"
	Yellow   = "\x1b[33m"
	Green    = "\x1b[32m"
	Rst      = "\x1b[0m"
	DarkGrey = "\x1b[90m"
)

var DefaultTimeLayout = "15:04:05"

func pad(s string, overallLen int) string {
	padStr := " "
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

type ColoredFormatter struct {
	// Sets the message time format
	TimeLayout string

	// If set to true then do not print extra fields
	ExcludeFields bool
}

func (f ColoredFormatter) Format(e *logrus.Entry) ([]byte, error) {
	buf := new(bytes.Buffer)

	var levelColor string
	switch e.Level {
	case logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel:
		levelColor = Red
	case logrus.WarnLevel:
		levelColor = Yellow
	case logrus.InfoLevel:
		levelColor = Green
	}
	tl := DefaultTimeLayout
	if f.TimeLayout != "" {
		tl = f.TimeLayout
	}

	level := fmt.Sprintf("[%s%s\x1b[0m]", levelColor, e.Level.String())

	fmt.Fprintf(buf, "[\x1b[90m%s\x1b[0m] %s  %s", e.Time.Format(tl), pad(level, 16), e.Message)

	if !f.ExcludeFields {
		flds := []string{}
		for k, v := range e.Data {
			if s, ok := v.(string); ok && k != "Level" && k != "Message" {
				flds = append(flds, DarkGrey+k+Rst+":"+DarkGrey+s+Rst)
			}
		}
		fmt.Fprintf(buf, " \x1b[36m%s\x1b[0m", strings.Join(flds, " "))
	}

	fmt.Fprint(buf, "\n")

	return buf.Bytes(), nil
}
