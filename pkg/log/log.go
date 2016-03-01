// log.go
package log

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

const (
	logDir = "/var/log/scloud"
)

var (
	logFileName string

	logFile *os.File
)

func InitLog() error {
	err := os.MkdirAll(logDir, 0644)

	if err != nil {
		return err
	}

	currentTime := time.Now().Unix()

	tm := time.Unix(currentTime, 0)

	timestamp := tm.Format("20060102150405")

	logBaseName := fmt.Sprintf("scloud_%s.log", timestamp)

	logFileName = filepath.Join(logDir, logBaseName)

	logFile, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)

	if err == nil {
		logFile.WriteString(fmt.Sprintf("%v Opened logfile at %v", os.Getpid(), time.Now()))
		os.Stderr = logFile
		syscall.Dup2(int(logFile.Fd()), 2)
	}
}

func CloseLog() {
	if logFile != nil {
		logFile.Close()
	}
}
