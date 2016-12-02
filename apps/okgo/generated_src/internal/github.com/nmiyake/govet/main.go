package amalgomated

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func AmalgomatedMain() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Failed to get working directory: %v\n", err)
		os.Exit(1)
	}

	// wd may be a symlink, so resolve to physical path
	physicalWd, err := filepath.EvalSymlinks(wd)
	if err != nil {
		fmt.Printf("Failed to evaluate symlinks for path %v: %v\n", wd, err)
		os.Exit(1)
	}

	cmd := exec.Command("go", append([]string{"vet"}, os.Args[1:]...)...)
	cmd.Env = append([]string{fmt.Sprintf("PWD=%v", physicalWd)}, os.Environ()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
