package consola

import "github.com/sirupsen/logrus"

func Example() {
	log := logrus.New()
	// Set the logrus formatter to a new consola.Formatter
	log.Formatter = Formatter{}

	log.Info("Hello, World!")
}
