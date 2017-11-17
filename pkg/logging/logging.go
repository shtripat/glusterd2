// Pacakge logging implements a common log intialization for GD2 and its CLI
package logging

import (
	"io"
	stdlog "log"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	DirFlag   = "logdir"
	DirHelp   = "Directory to store log files"
	FileFlag  = "logfile"
	FileHelp  = "Name for log file"
	LevelFlag = "loglevel"
	LevelHelp = "Severity of messages to be logged"
)

var logWriter io.WriteCloser

func openLogFile(filepath string) (io.WriteCloser, error) {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func setLogOutput(w io.Writer) {
	log.SetOutput(w)
	stdlog.SetOutput(log.StandardLogger().Writer())
}

// Init initializes the default logrus logger
// Should be called as early as possible when a process starts.
// Note that this does not create a new logger. Packages should still continue
// importing and using logrus as before.
func Init(logdir string, logFileName string, logLevel string) error {
	// Close the previously opened Log file
	if logWriter != nil {
		logWriter.Close()
		logWriter = nil
	}

	l, err := log.ParseLevel(strings.ToLower(logLevel))
	if err != nil {
		setLogOutput(os.Stderr)
		log.WithError(err).Debug("Failed to parse log level")
		return err
	}
	log.SetLevel(l)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	if strings.ToLower(logFileName) == "stderr" || logFileName == "-" {
		setLogOutput(os.Stderr)
	} else if strings.ToLower(logFileName) == "stdout" {
		setLogOutput(os.Stdout)
	} else {
		logFilePath := path.Join(logdir, logFileName)
		logFile, err := openLogFile(logFilePath)
		if err != nil {
			setLogOutput(os.Stderr)
			log.WithError(err).Debug("Failed to open log file %s", logFilePath)
			return err
		}
		setLogOutput(logFile)
		logWriter = logFile
	}
	return nil
}
