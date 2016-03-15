package engine

import (
	"errors"
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"strings"
	"syscall"
	"time"
)

const (
	DefaultTimeOut time.Duration = 60 * time.Minute
)

var (
	CommandNotStartError = errors.New("The command is not started.")

	CommandStillRunning = errors.New("The command is still running.")
)

type Command struct {
	Name string
	Args []string
	Pwd  string
	Env
	Stdout  io.Writer
	stderr  io.Writer
	Timeout time.Duration
	status  Status
	end     time.Time
	start   time.Time
}

type Status int

const (
	StatusOK       Status = 0
	StatusErr      Status = 1
	StatusNotFound Status = 127
)

type Env map[string]string

func (this *Command) PrintCommand() string {
	return fmt.Sprintf("==> Executing: %s %s\n", this.Name, strings.Join(this.Args, " "))
}

func (this *Command) Run() error {

	this.start = time.Now()

	_, err := osexec.LookPath(name)

	if err != nil {
		this.end = time.Now()
		this.status = StatusNotFound
		return err
	}

	cmd := osexec.Command(this.Name, this.Args)

	this.setCurrentEnv(cmd)

	// set output
	if this.Stdout != nil {
		cmd.Stdout = this.Stdout
	}

	// set stderr
	if this.stderr != nil {
		cmd.Stderr = this.stderr
	}

	cmd.Stdout.Write([]byte(this.PrintCommand()))

	startErr := cmd.Start()

	if startErr != nil {

		this.end = time.Now()
		this.status = StatusErr
		return startErr
	}

	timer := time.NewTimer(this.Timeout)

	go func(time *time.Timer, cmd *osexec.Cmd) {
		for _ = range timer.C {
			err := cmd.Process.Signal(os.Kill)

			_, err := cmd.Stderr.Write([]byte(err.Error()))

			fmt.Fprintf(os.Stderr, err)
			fmt.Fprintf(os.Stdout, err)
		}
	}(timer, cmd)

	err := cmd.Wait()

	if err != nil {
		fmt.Fprintf(os.Stderr, err)
		fmt.Fprintf(os.Stdout, err)
		return err
	}

	return nil
}

// set the env for the current command
func (this *Command) setCurrentEnv(cmd *osexec.Cmd) []string {

	env := os.Environ()

	if this.Env != nil {
		for k, v := range this.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	return env
}

// Get the status of the command
func (this *Command) StatusCode() (int, error) {

	if this.start.IsZero() {
		return -1, CommandNotStartError
	}

	if this.end.IsZero() {
		return -1, CommandStillRunning
	}

	return int(this.status)
}
