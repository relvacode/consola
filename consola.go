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
var DefaultFieldSeparator = ":"

type ColoredFormatter struct {
	// Sets the message time format
	TimeLayout string

	// If set to true then do not print extra fields
	ExcludeFields bool

	// String value used to separate log fields
	FieldSeparator string
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

	fsep := DefaultFieldSeparator
	if f.FieldSeparator != "" {
		fsep = f.FieldSeparator
	}

	level := fmt.Sprintf("[%s%s\x1b[0m]", levelColor, e.Level.String())

	fmt.Fprintf(buf, "[\x1b[90m%s\x1b[0m] %s  %s", e.Time.Format(tl), level, e.Message)

	if !f.ExcludeFields {
		flds := []string{}
		for k, v := range e.Data {
			if s, ok := v.(string); ok && k != "Level" && k != "Message" {
				flds = append(flds, DarkGrey+k+Rst+fsep+DarkGrey+s+Rst)
			}

		}
		fmt.Fprintf(buf, " \x1b[36m%s\x1b[0m", strings.Join(flds, " "))
	}

	fmt.Fprint(buf, "\n")

	return buf.Bytes(), nil
}
