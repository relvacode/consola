// consola is a simple package for printing better looking Logrus log messages to the console.
package consola

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

// Terminal color constants
const (
	Red      = "\x1b[31m"
	Yellow   = "\x1b[33m"
	Green    = "\x1b[32m"
	Rst      = "\x1b[0m"
	DarkGrey = "\x1b[90m"
)

var DefaultTimeLayout = "15:04:05 02/01"
var DefaultFieldSeparator = ":"

type level struct {
	logrus.Level
}

func (l level) Write(buf io.Writer, color bool) {
	var str = fmt.Sprintf("%-7s", l.String())
	if color {
		var clr string
		switch l.Level {
		case logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel:
			clr = Red
		case logrus.WarnLevel:
			clr = Yellow
		case logrus.InfoLevel:
			clr = Green
		}
		str = clr + str + Rst
	}
	fmt.Fprint(buf, str, " ")
}

type fields struct {
	logrus.Fields
}

func (f fields) WritePrefix(buf io.Writer, color bool, prefix string) {
	if pk, ok := f.Fields[prefix]; ok {
		var str string = fmt.Sprint(pk)
		if color {
			str = DarkGrey + str + Rst
		}
		fmt.Fprint(buf, str, "\t")
		delete(f.Fields, prefix)
	}
}

func (f fields) Write(buf io.Writer, color bool, sep string) {
	l := len(f.Fields)
	var n int
	for k, v := range f.Fields {
		n++
		if s, ok := v.(string); ok && k != "Level" && k != "Message" {
			if color {
				fmt.Fprint(buf, DarkGrey+k+Rst+sep+DarkGrey+s+Rst)
			} else {
				fmt.Fprint(buf, k+sep+s)
			}
			if l != n {
				fmt.Fprint(buf, " ")
			}
		}
	}
}

type Formatter struct {
	// Sets the message time format
	TimeLayout string

	// If set to true then do not print extra fields
	ExcludeFields bool

	// String value used to separate log fields
	FieldSeparator string

	PrimaryPrefixField string

	SecondaryPrefixField string

	// Enable color
	Color bool
}

func (f Formatter) Format(e *logrus.Entry) ([]byte, error) {
	buf := new(bytes.Buffer)

	layout := DefaultTimeLayout
	if f.TimeLayout != "" {
		layout = f.TimeLayout
	}

	sep := DefaultFieldSeparator
	if f.FieldSeparator != "" {
		sep = f.FieldSeparator
	}

	fields := fields{e.Data}

	fmt.Fprint(buf, e.Time.Format(layout), " ")
	level{e.Level}.Write(buf, f.Color)

	if f.PrimaryPrefixField != "" {
		fields.WritePrefix(buf, f.Color, f.PrimaryPrefixField)
	}
	if f.SecondaryPrefixField != "" {
		fields.WritePrefix(buf, f.Color, f.SecondaryPrefixField)
	}

	fmt.Fprint(buf, e.Message, " ")

	if !f.ExcludeFields {
		fields.Write(buf, f.Color, sep)
	}

	fmt.Fprint(buf, "\n")

	return buf.Bytes(), nil
}
