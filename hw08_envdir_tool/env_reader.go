package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading dir: %s, %w", dir, err)
	}
	env := make(Environment)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.Contains(entry.Name(), "=") {
			return nil, fmt.Errorf("incorrect file name: %s", entry.Name())
		}
		info, err := entry.Info()
		if err != nil {
			return nil, fmt.Errorf("error getting file info, %w", err)
		}
		if info.Size() == 0 {
			env[info.Name()] = EnvValue{NeedRemove: true}
			continue
		}
		line, err := readEnvFile(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("error reading file content, %w", err)
		}
		env[info.Name()] = EnvValue{Value: *line}
	}
	return env, nil
}

func readEnvFile(filePath string) (*string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var line string
	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line = scanner.Text()

		line = strings.ReplaceAll(line, "\x00", "\n")
		line = strings.TrimRight(line, " \t")
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &line, nil
}
