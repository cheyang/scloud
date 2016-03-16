package engine_test

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	. "github.com/cheyang/scloud/pkg/engine"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test command", func() {

	Context("#RUn", func() {

		It("output to log file", func() {
			cmd := NewCommand("ls", "-l")
			cmd.SetWorkingDir("/tmp")

			currentTime := time.Now().Unix()

			tm := time.Unix(currentTime, 0)

			timestamp := tm.Format("20060102150405")

			logBaseName := fmt.Sprintf("test_%s.log", timestamp)

			logFileName := filepath.Join("/tmp", logBaseName)

			logFile, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)

			if err != nil {
				fmt.Println("error:", err)
			}

			logFile.WriteString(fmt.Sprintf("%v Opened logfile at %v", os.Getpid(), time.Now()))

			cmd.SetStdout(logFile)

			cmd.SetStderr(logFile)

			err := cmd.Run()

			if err != nil {
				fmt.Println("error:", err)
			}

			Expect(err).To(BeNil())

		})

	})

})
