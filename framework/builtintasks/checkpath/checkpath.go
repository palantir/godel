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
	"io"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/termie/go-shutil"
)

const (
	CmdName       = "check-path"
	ApplyFlagName = "apply"
)

func VerifyProject(wd string, apply bool, stdout io.Writer) error {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return errors.Errorf("GOPATH environment variable must be set")
	}

	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		// GOROOT is not set as environment variable: check value of GOROOT provided by compiled code
		compiledGoRoot := runtime.GOROOT()
		if goRootPathInfo, err := os.Stat(compiledGoRoot); err != nil || !goRootPathInfo.IsDir() {
			// compiled GOROOT does not exist: get GOROOT from "go env GOROOT"
			goEnvRoot := goEnvRoot()

			fmt.Fprintln(stdout, "GOROOT environment variable is not set and the value provided in the compiled code does not exist locally.")

			if goEnvRoot != "" {
				shellCfgFiles := getShellCfgFiles()
				if !apply {
					fmt.Fprintf(stdout, "'go env' reports that GOROOT is %s. Suggested fix:\n", goEnvRoot)
					fmt.Fprintf(stdout, "  export GOROOT=%s\n", goEnvRoot)
					for _, currCfgFile := range shellCfgFiles {
						fmt.Fprintf(stdout, "  echo \"export GOROOT=%s\" >> %q \n", goEnvRoot, currCfgFile)
					}
				} else {
					for _, currCfgFile := range shellCfgFiles {
						fmt.Fprintf(stdout, "Adding \"export GOROOT=%s\" to %s...\n", goEnvRoot, currCfgFile)
						if err := appendToFile(currCfgFile, fmt.Sprintf("export GOROOT=%v\n", goEnvRoot)); err != nil {
							fmt.Fprintf(stdout, "Failed to add \"export GOROOT=%s\" in %s\n", goEnvRoot, currCfgFile)
						}
					}
				}
			} else {
				fmt.Fprintln(stdout, "Unable to determine GOROOT using 'go env GOROOT'. Ensure that Go was installed properly with source files.")
			}
		}
	}

	srcPath, err := gitRepoRootPath(wd)
	if err != nil {
		return errors.Errorf("Directory %q must be in a git project", wd)
	}

	pkgPath, gitRemoteURL, err := gitRemoteOriginPkgPath(wd)
	if err != nil {
		return errors.Errorf("Unable to determine URL of git remote for %s: verify that the project was checked out from a git remote and that the remote URL is set", wd)
	} else if pkgPath == "" {
		return errors.Errorf("Unable to determine the expected package path from git remote URL of %s", gitRemoteURL)
	}

	dstPath := path.Join(gopath, "src", pkgPath)
	if srcPath == dstPath {
		fmt.Fprintln(stdout, "Project appears to be in the correct location")
		return nil
	}
	dstPathParentDir := path.Dir(dstPath)

	fmt.Fprintf(stdout, "Project path %q differs from expected path %q\n", srcPath, dstPath)
	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		fmt.Fprintf(stdout, "Expected destination path %q already exists.\n", dstPath)
		fmt.Fprintln(stdout, "If this project is known to be the correct one, remove or rename the file or directory at the destination path and run this command again.")
		return nil
	}

	var fixOperation func(srcPath, dstPath string) error
	var fixOpMessage string
	if pathsOnSameDevice(srcPath, gopath) {
		if !apply {
			fmt.Fprintln(stdout, "Project and GOPATH are on same device. Suggested fix:")
			fmt.Fprintf(stdout, "  mkdir -p %q && mv %q %q\n", dstPathParentDir, srcPath, dstPath)
		}
		fixOpMessage = fmt.Sprintf("Moving %q to %q...", wd, dstPath)
		fixOperation = os.Rename
	} else {
		if !apply {
			fmt.Fprintln(stdout, "Project and GOPATH are on different devices. Suggested fix:")
			fmt.Fprintf(stdout, "  mkdir -p %q && cp -r %q %q\n", dstPathParentDir, srcPath, dstPath)
		}
		fixOpMessage = fmt.Sprintf("Copying %q to %q...", wd, dstPath)
		fixOperation = func(srcPath, dstPath string) error {
			return shutil.CopyTree(srcPath, dstPath, nil)
		}
	}

	if !apply {
		fmt.Fprintf(stdout, "%s found issues. Run the following godel command to implement suggested fixes:\n", CmdName)
		fmt.Fprintf(stdout, "  %s --%s\n", CmdName, ApplyFlagName)
	} else {
		if err := os.MkdirAll(dstPathParentDir, 0755); err != nil {
			return errors.Errorf("Failed to create path to %q: %v", dstPathParentDir, err)
		}
		fmt.Fprintln(stdout, fixOpMessage)
		if err := fixOperation(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
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
