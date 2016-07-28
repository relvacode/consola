package consola

import "github.com/Sirupsen/logrus"

func Example() {
	log := logrus.New()
	// Set the logrus formatter to a new consola.ColoredFormatter
	log.Formatter = ColoredFormatter{}

	log.Info("Hello, World!")
}
