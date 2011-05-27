package log

import (
	"io"
	stdlog "log"
	"os"
)

var (
	Error  *stdlog.Logger
	Warn   *stdlog.Logger
	Info   *stdlog.Logger
	Debug  *stdlog.Logger
	output io.Writer
)

func init() {
	SetLog(os.Stderr)
}

func SetFile(filename string) os.Error {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	SetLog(logfile)
	return nil
}

func SetLog(logfile io.Writer) {
	output = logfile
	Error = stdlog.New(logfile, "[E] ", stdlog.Ldate|stdlog.Lmicroseconds)
	Warn = stdlog.New(logfile, "[W] ", stdlog.Ldate|stdlog.Lmicroseconds)
	Info = stdlog.New(logfile, "[I] ", stdlog.Ldate|stdlog.Lmicroseconds)
	Debug = stdlog.New(logfile, "[D] ", stdlog.Ldate|stdlog.Lmicroseconds)
}

// After setting the log to go to a file, this can be used to also show the
// log in the console.  Calling this when the log is already being printed to
// standard error will cause duplicates to be written.
func ShowInConsole() {
	SetLog(io.MultiWriter(output, os.Stderr))
}
