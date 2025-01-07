package logging

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type BagmanLogFormatter struct{}

func (f *BagmanLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	fmt.Fprintf(
		b, "<%s> - LogLevel: %s UUID: %s Message: %s\n",
		entry.Time.Format(time.RFC3339),
		entry.Level,
		entry.Data["UUID"],
		entry.Message,
	)
	return b.Bytes(), nil
}

func NewLogger(mode string, filename string) (*logrus.Logger, error) {
	var output *os.File

	switch mode {
	case "stdout":
		output = os.Stdout
	case "file":
		var err error
		output, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
	default:
		return nil, fmt.Errorf("unknown logger mode: %s, use 'file', or 'stdout'", mode)
	}

	logger := &logrus.Logger{
		Out:       output,
		Formatter: new(BagmanLogFormatter),
		Level:     logrus.InfoLevel,
	}

	return logger, nil
}
