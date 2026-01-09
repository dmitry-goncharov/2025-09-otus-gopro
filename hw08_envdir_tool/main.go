package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("You should use go-envdir with args: /path/to/env/dir command arguments")
		os.Exit(1)
	}
	env, err := ReadDir(args[1])
	if err != nil {
		fmt.Println("Error reading dir", args[1], err.Error())
		os.Exit(1)
	}
	resCode := RunCmd(args[2:], env)
	os.Exit(resCode)
}
