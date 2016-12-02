package main

import (
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("go", append([]string{"test"}, os.Args[1:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
