package main

import (
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
	return runCmd(cmd)
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

func runCmd(cmd []string) int {
	c := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	c.Env = os.Environ()
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		fmt.Println("Error run command", err.Error())
	}
	return c.ProcessState.ExitCode()
}
