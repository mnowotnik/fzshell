package output

import (
	"fmt"
	logging "log"
	"os"
	"path/filepath"
	"time"
)

func printErr(msg string) {
	fmt.Fprintln(os.Stderr, "fzshell: "+msg)
}

var logger *logging.Logger = nil

func Log() *logging.Logger {
	if logger != nil {
		return logger
	}
	terminalLogger := func() *logging.Logger {
		logger = logging.New(os.Stdout, "[fzshell]  ", logging.LstdFlags)
		return logger
	}
	if os.Getenv("FZSHELL_DEBUG") == "" {
		return terminalLogger()
	}
	var logF *os.File
	if os.Getenv("XDG_RUNTIME_DIR") != "" {
		dirPath := filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), "fzshell")
		err := os.Mkdir(dirPath, 0700)
		if err != nil && !os.IsExist(err) {
			printErr("Could not create directory for logs: " + dirPath)
			return terminalLogger()
		}
		logPath := filepath.Join(dirPath, "fzshell.log")
		logF, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			printErr("Could not create log file: " + logPath)
			return terminalLogger()
		}
		defer logF.Close()
	} else {
		var err error
		logPath := filepath.Join(os.TempDir(), "fzshell"+time.Now().Format("060102")+".log")
		logF, err = os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil && !os.IsExist(err) {
			printErr("Could not create log file in temp directory")
			return terminalLogger()
		}
		defer logF.Close()
		fi, err := os.Stat(logPath)
		if err != nil || fi.Mode().Perm() != 0600 {
			printErr("Could not create log file in temp directory")
			return terminalLogger()
		}
	}
	logger = logging.New(logF, "", logging.LstdFlags)
	return logger
}
