package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

type ExecOptions struct {
	Command string
	Dir     string
	Env     []string
}

type ExecResult struct {
	Stdout     string
	Stderr     string
	ExitStatus int
}

func Exec(eo ExecOptions) (ExecResult, error) {
	var err error
	var cmd *exec.Cmd
	var er ExecResult

	var stdout io.ReadCloser
	var stderr io.ReadCloser
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	x := strings.Split(eo.Command, " ")
	if len(x) > 1 {
		cmd = exec.Command(x[0], x[1:]...)
	} else {
		cmd = exec.Command(x[0])
	}

	if len(eo.Env) > 0 {
		cmd.Env = eo.Env
	}

	if eo.Dir != "" {
		cmd.Dir = eo.Dir
	}

	stdout, err = cmd.StdoutPipe()
	if err != nil {
		return er, err
	}

	stderr, err = cmd.StderrPipe()
	if err != nil {
		return er, err
	}

	if err := cmd.Start(); err != nil {
		return er, err
	}

	if _, err := bufio.NewReader(stderr).WriteTo(&stderrBuf); err != nil {
		return er, err
	}

	if _, err := bufio.NewReader(stdout).WriteTo(&stdoutBuf); err != nil {
		return er, err
	}

	if err := cmd.Wait(); err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			er.ExitStatus = int(exit.ProcessState.Sys().(syscall.WaitStatus) / 256)
		}
	}

	er.Stdout = stdoutBuf.String()
	er.Stderr = stderrBuf.String()

	return er, err
}

func RequiredCommands(commands []string) error {
	for _, command := range commands {
		_, err := exec.LookPath(command)
		if err != nil {
			return fmt.Errorf("Required command not found: %s", command)
		}
	}

	return nil
}
