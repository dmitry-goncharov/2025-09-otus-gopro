package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		panic("no command")
	}
	err := updateEnv(env)
	if err != nil {
		fmt.Println("Error update environment", err.Error())
		return 1
	}
	err = runCmd(cmd)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		fmt.Println("Error run command", err.Error())
		return 1
	}
	fmt.Println("Run command successfully", cmd)
	return 0
}

func updateEnv(env Environment) error {
	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return err
			}
		} else {
			err := os.Setenv(key, val.Value)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func runCmd(cmd []string) error {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return err
	}
	return nil
}
