package cliutil

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/pkg/errors"
)

func Exec(command string, args []string) (string, string, int, error) {
	if len(command) < 1 {
		return "", "", 0, errors.New("invalid command")
	}

	return ExecWithContext(context.TODO(), command, args)
}

func ShellExec(command string) (string, string, int, error) {
	if len(command) < 1 {
		return "", "", 0, errors.New("invalid command")
	}
	args := []string{"-c", command}
	return ExecWithContext(context.TODO(), "/bin/sh", args)
}

func ExecWithContext(ctx context.Context, command string, args []string) (string, string, int, error) {
	var (
		err      error
		stderr   bytes.Buffer
		stdout   bytes.Buffer
		code     int
		output   string
		warnings string
	)

	if ctx == nil {
		ctx = context.TODO()
	}
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	output = stdout.String()
	warnings = stderr.String()
	code = cmd.ProcessState.ExitCode()

	if code != 0 && err == nil {
		err = errors.Errorf("execute code %d: `%s`", code, warnings)
	}

	return output, warnings, code, err
}
