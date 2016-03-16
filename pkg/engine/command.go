package engine

import (
	"errors"
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"strings"

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
	Stdout   io.Writer
	Stderr   io.Writer
	Timeout  time.Duration
	status   Status
	end      time.Time
	start    time.Time
	duration time.Duration
}

type Status int

const (
	StatusOK       Status = 0
	StatusErr      Status = 1
	StatusNotFound Status = 127
)

type Env map[string]string

func NewCommand(name string, args ...string) *Command {
	return &Command{
		Name:    name,
		Args:    append([]string{}, args...),
		Timeout: DefaultTimeOut,
		Env:     Env(make(map[string]string)),
	}
}

// Set environment variables
func (this *Command) SetEnvironmentVar(key string, value string) {
	this.Env[key] = value
}

// Set working dir
func (this *Command) SetWorkingDir(dir string) {
	this.Pwd = dir
}

func (this *Command) SetStdout(output io.Writer) {
	this.Stdout = output
}

func (this *Command) SetStderr(stderr io.Writer) {
	this.Stderr = stderr
}

func (this *Command) PrintCommand() string {
	return fmt.Sprintf("==> Executing: %s %s\n", this.Name, strings.Join(this.Args, " "))
}

// Run executes the command and blocks until the command completes.
// If the command returns a failure status, an error is returned
// which includes the status.
func (this *Command) Run() error {

	this.start = time.Now()

	defer func() {
		this.end = time.Now()
		this.duration = this.end.Sub(this.start)
	}()

	_, err := osexec.LookPath(this.Name)

	if err != nil {
		this.status = StatusNotFound
		return err
	}

	cmd := osexec.Command(this.Name, this.Args...)

	cmd.Dir = this.Pwd

	this.setCurrentEnv(cmd)

	// set output
	if this.Stdout != nil {
		cmd.Stdout = this.Stdout
	}

	// set stderr
	if this.Stderr != nil {
		cmd.Stderr = this.Stderr
	}

	cmd.Stdout.Write([]byte(this.PrintCommand()))

	startErr := cmd.Start()

	if startErr != nil {
		this.status = StatusErr
		return startErr
	}

	timer := time.NewTimer(this.Timeout)

	go func(time *time.Timer, cmd *osexec.Cmd) {
		for _ = range timer.C {
			err := cmd.Process.Signal(os.Kill)

			_, err = cmd.Stderr.Write([]byte(err.Error()))

			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stdout, err)
		}
	}(timer, cmd)

	err = cmd.Wait()

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stdout, err)
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

	return int(this.status), nil
}

func (this *Command) GetPeriod() time.Duration {
	return this.duration
}
