// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package checkpath

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"github.com/termie/go-shutil"
)

func VerifyProject(wd string, info bool) error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("GOPATH environment variable must be set")
	}

	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		// GOROOT is not set as environment variable: check value of GOROOT provided by compiled code
		compiledGoRoot := runtime.GOROOT()
		if goRootPathInfo, err := os.Stat(compiledGoRoot); err != nil || !goRootPathInfo.IsDir() {
			// compiled GOROOT does not exist: get GOROOT from "go env GOROOT"
			goEnvRoot := goEnvRoot()

			if info {
				fmt.Printf("GOROOT environment variable is not set and the value provided in the compiled code does not exist locally.\n")
			}

			if goEnvRoot != "" {
				shellCfgFiles := getShellCfgFiles()
				if info {
					fmt.Printf("'go env' reports that GOROOT is %v. Suggested fix:\n\texport GOROOT=%v\n", goEnvRoot, goEnvRoot)
					for _, currCfgFile := range shellCfgFiles {
						fmt.Printf("\techo \"export GOROOT=%v\" >> %q \n", goEnvRoot, currCfgFile)
					}
				} else {
					for _, currCfgFile := range shellCfgFiles {
						fmt.Printf("Adding \"export GOROOT=%v\" to %v...\n", goEnvRoot, currCfgFile)
						if err := appendToFile(currCfgFile, fmt.Sprintf("export GOROOT=%v\n", goEnvRoot)); err != nil {
							fmt.Printf("Failed to add \"export GOROOT=%v\" in %v\n", goEnvRoot, currCfgFile)
						}
					}
				}
			} else {
				fmt.Printf("Unable to determine GOROOT using 'go env GOROOT'. Ensure that Go was installed properly with source files.\n")
			}
		}
	}

	srcPath, err := gitRepoRootPath(wd)
	if err != nil {
		return fmt.Errorf("Directory %q must be in a git project", wd)
	}

	pkgPath, gitRemoteURL, err := gitRemoteOriginPkgPath(wd)
	if err != nil {
		return fmt.Errorf("Unable to determine URL of git remote for %v: verify that the project was checked out from a git remote and that the remote URL is set", wd)
	} else if pkgPath == "" {
		return fmt.Errorf("Unable to determine the expected package path from git remote URL of %v", gitRemoteURL)
	}

	dstPath := path.Join(gopath, "src", pkgPath)
	if srcPath == dstPath {
		fmt.Printf("Project appears to be in the correct location\n")
		return nil
	}
	dstPathParentDir := path.Dir(dstPath)

	fmt.Printf("Project path %q differs from expected path %q\n", srcPath, dstPath)
	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		if !info {
			return fmt.Errorf("Destination path %q already exists", dstPath)
		}
		fmt.Printf("Expected destination path %q already exists.\nIf this project is known to be the correct one, remove or rename the file or directory at the destination path and run this command again.\n", dstPath)
		return nil
	}

	var fixOperation func(srcPath, dstPath string) error
	var fixOpMessage string
	if pathsOnSameDevice(srcPath, gopath) {
		if info {
			fmt.Printf("Project and GOPATH are on same device. Suggested fix:\n\tmkdir -p %q && mv %q %q\n", dstPathParentDir, srcPath, dstPath)
		}
		fixOpMessage = fmt.Sprintf("Moving %q to %q...\n", wd, dstPath)
		fixOperation = os.Rename
	} else {
		if info {
			fmt.Printf("Project and GOPATH are on different devices. Suggested fix:\n\tmkdir -p %q && cp -r %q %q\n", dstPathParentDir, srcPath, dstPath)
		}
		fixOpMessage = fmt.Sprintf("Copying %q to %q...\n", wd, dstPath)
		fixOperation = func(srcPath, dstPath string) error {
			return shutil.CopyTree(srcPath, dstPath, nil)
		}
	}

	if info {
		fmt.Printf("\n%v found issues. Run the following godel command to implement suggested fixes:\n\t%v\n", cmd, cmd)
	} else {
		if err := os.MkdirAll(dstPathParentDir, 0755); err != nil {
			return fmt.Errorf("Failed to create path to %q: %v", dstPathParentDir, err)
		}
		fmt.Printf(fixOpMessage)
		if err := fixOperation(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func VerifyGoEnv(wd string) {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		// GOROOT is not set as environment variable: check value of GOROOT provided by compiled code
		compiledGoRoot := runtime.GOROOT()
		if info, err := os.Stat(compiledGoRoot); err != nil || !info.IsDir() {
			// compiled GOROOT does not exist: get GOROOT from "go env GOROOT" command and display warning
			goEnvRoot := goEnvRoot()

			title := "GOROOT environment variable is empty"
			content := fmt.Sprintf("GOROOT is required to build Go code. The GOROOT environment variable was not set and the value of GOROOT in the compiled code does not exist locally.\n")
			content += fmt.Sprintf("Run the godel command '%v --info' for more information.", cmd)

			if goEnvRoot != "" {
				content += fmt.Sprintf("\nFalling back to value provided by 'go env GOROOT': %v", goEnvRoot)
				if err := os.Setenv("GOROOT", goEnvRoot); err != nil {
					fmt.Printf("Failed to set GOROOT environment variable: %v", err)
					return
				}
			}
			printWarning(title, content)
		}
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		title := "GOPATH environment variable is empty"
		content := fmt.Sprintf("GOPATH is required to build Go code.\nSee https://golang.org/doc/code.html#GOPATH for more information.")
		printWarning(title, content)
		return
	}

	goSrcPath := path.Join(gopath, "src")
	if !isSubdir(goSrcPath, wd) {
		title := "project directory is not a subdirectory of $GOPATH/src"
		content := fmt.Sprintf("%q is not a subdirectory of the GOPATH/src directory (%q): the project location is likely incorrect.\n", wd, goSrcPath)
		content += fmt.Sprintf("Run the godel command '%v --info' for more information.", cmd)
		printWarning(title, content)
	}
}

func goEnvRoot() string {
	if output, err := exec.Command("go", "env", "GOROOT").CombinedOutput(); err == nil {
		return strings.TrimSpace(string(output))
	}
	return ""
}

func getShellCfgFiles() []string {
	var files []string
	for _, currCfg := range []string{".bash_profile", ".zshrc"} {
		if currFile := checkForCfgFile(currCfg); currFile != "" {
			files = append(files, currFile)
		}
	}
	return files
}

func appendToFile(filename, content string) (rErr error) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close file %s in defer", filename)
		}
	}()
	if _, err = f.WriteString(content); err != nil {
		return err
	}
	return nil
}

func checkForCfgFile(name string) string {
	cfgFile := path.Join(os.Getenv("HOME"), name)
	if _, err := os.Stat(cfgFile); err == nil {
		return cfgFile
	}
	return ""
}

func printWarning(title, content string) {
	warning := fmt.Sprintf("| WARNING: %v |", title)
	fmt.Println(strings.Repeat("-", len(warning)))
	fmt.Println(warning)
	fmt.Println(strings.Repeat("-", len(warning)))
	fmt.Println(content)
	fmt.Println()
}

func pathsOnSameDevice(p1, p2 string) bool {
	id1, ok := getDeviceID(p1)
	if !ok {
		return false
	}
	id2, ok := getDeviceID(p2)
	if !ok {
		return false
	}
	return id1 == id2
}

func getDeviceID(p string) (interface{}, bool) {
	fi, err := os.Stat(p)
	if err != nil {
		return 0, false
	}
	s := fi.Sys()
	switch s := s.(type) {
	case *syscall.Stat_t:
		return s.Dev, true
	}
	return 0, false
}

func isSubdir(base, dst string) bool {
	relPath, err := filepath.Rel(base, dst)
	return err == nil && !strings.HasPrefix(relPath, "..")
}

func gitRepoRootPath(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%v failed: could not determine root of git repository for directory %v", cmd.Args, dir)
	}
	return strings.TrimSpace(string(output)), nil
}

func gitRemoteOriginPkgPath(dir string) (string, string, error) {
	cmd := exec.Command("git", "ls-remote", "--get-url")
	cmd.Dir = dir
	bytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("%v failed: git remote URL is not set", cmd.Args)
	}

	pkgPath := ""
	url := strings.TrimSpace(string(bytes))
	if protocolSeparator := regexp.MustCompile("https?://").FindStringIndex(url); protocolSeparator != nil {
		if protocolSeparator[0] == 0 {
			pkgPath = url[protocolSeparator[1]:]
		}
	} else if atSign := strings.Index(url, "@"); atSign != -1 {
		// assume SSH format of "user@host:org/repo"
		pkgPath = url[atSign+1:]
		pkgPath = strings.Replace(pkgPath, ":", "/", 1)
	}

	// trim ".git" suffix if present
	if strings.HasSuffix(pkgPath, ".git") {
		pkgPath = pkgPath[:len(pkgPath)-len(".git")]
	}

	return pkgPath, url, nil
}
