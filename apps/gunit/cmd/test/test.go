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

package test

import (
	"bufio"
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/gunit/cmd"
	"github.com/palantir/godel/apps/gunit/config"
)

const (
	goTestCmdName          = "test"
	goCoverCmdName         = "cover"
	coverageOutputPathFlag = "coverage-output"
	pkgsParamName          = "packages"
)

var (
	goTestCmd        = cmd.Library.MustNewCmd("gotest")
	goCoverCmd       = cmd.Library.MustNewCmd("gotest")
	goJUnitReportCmd = cmd.Library.MustNewCmd("gojunitreport")
	gtCmd            = cmd.Library.MustNewCmd("gt")
)

var goTestCmdGenerator = &testCmdGenerator{
	cmd:          goTestCmd,
	usage:        "Runs 'go test' on project packages",
	runFunc:      runGoTest,
	paramCreator: testParamCreator,
}

func GoTestCommand(supplier amalgomated.CmderSupplier) cli.Command {
	cmd := goTestCmdGenerator.baseTestCmd(supplier)
	cmd.Name = goTestCmdName
	return cmd
}

func RunGoTestAction(supplier amalgomated.CmderSupplier) func(ctx cli.Context) error {
	return func(ctx cli.Context) error {
		testParams := testParamCreator(ctx)
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return err
		}
		return goTestCmdGenerator.runTestCmdForPkgs(nil, testParams, supplier, wd, ctx.App.OnExit)
	}
}

func GTCommand(supplier amalgomated.CmderSupplier) cli.Command {
	gtCmd := &testCmdGenerator{
		cmd:          gtCmd,
		usage:        "Runs 'gt' on project packages",
		runFunc:      runGoTest,
		paramCreator: testParamCreator,
	}
	return gtCmd.baseTestCmd(supplier)
}

func GoCoverCommand(supplier amalgomated.CmderSupplier) cli.Command {
	coverCmd := &testCmdGenerator{
		cmd:          goCoverCmd,
		usage:        "Runs 'go cover' on project packages",
		runFunc:      runGoTestCover,
		paramCreator: coverageParamCreator,
	}
	cmd := coverCmd.baseTestCmd(supplier)
	cmd.Name = goCoverCmdName
	cmd.Flags = append(cmd.Flags,
		flag.StringFlag{
			Name:  coverageOutputPathFlag,
			Usage: "Path to coverage output file",
		},
	)
	return cmd
}

type testCmdGenerator struct {
	cmd          amalgomated.Cmd
	usage        string
	runFunc      runTestFunc
	paramCreator paramCreatorFunc
}

func (g *testCmdGenerator) baseTestCmd(supplier amalgomated.CmderSupplier) cli.Command {
	return cli.Command{
		Name:  g.cmd.Name(),
		Usage: g.usage,
		Flags: []flag.Flag{
			flag.StringSlice{
				Name:     pkgsParamName,
				Usage:    "Packages to test",
				Optional: true,
			},
		},
		Action: func(ctx cli.Context) error {
			pkgsParam := ctx.Slice(pkgsParamName)
			testParams := g.paramCreator(ctx)
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return g.runTestCmdForPkgs(pkgsParam, testParams, supplier, wd, ctx.App.OnExit)
		},
	}
}

func (g *testCmdGenerator) runTestCmdForPkgs(pkgsParam []string, testParams testCtxParams, supplier amalgomated.CmderSupplier, wd string, onExitManager cli.OnExit) error {
	cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	// tagsMatcher is a matcher that matches the specified tags (or nil if no tags were specified)
	tagsMatcher, err := cmd.TagsMatcher(testParams.tags, cfg)
	if err != nil {
		return err
	}
	excludeMatchers := []matcher.Matcher{cfg.Exclude}
	if tagsMatcher != nil {
		// if tagsMatcher is non-nil, should exclude all files that do not match the tags
		excludeMatchers = append(excludeMatchers, matcher.Not(tagsMatcher))
	}

	pkgs, err := cmd.PkgPaths(pkgsParam, wd, matcher.Any(excludeMatchers...))
	if err != nil {
		return err
	}

	if len(pkgs) == 0 {
		return errors.Errorf("no packages to test")
	}

	placeholderFiles, err := createPlaceholderTestFiles(pkgs, wd)
	if err != nil {
		return errors.Wrapf(err, "failed to create placeholder files for packages %v", pkgs)
	}

	cleanup := func() {
		for _, currFileToRemove := range placeholderFiles {
			if err := os.Remove(currFileToRemove); err != nil {
				fmt.Printf("%+v\n", errors.Wrapf(err, "failed to remove file %s", currFileToRemove))
			}
		}
	}

	// register cleanup task on exit
	cleanupID := onExitManager.Register(cleanup)

	// clean up placeholder files after function
	defer func() {
		// unregister cleanup task from CLI
		onExitManager.Unregister(cleanupID)

		// run cleanup
		cleanup()
	}()

	return g.runTestCmd(supplier, pkgs, testParams, wd)
}

func (g *testCmdGenerator) runTestCmd(supplier amalgomated.CmderSupplier, pkgs []string, params testCtxParams, wd string) (rErr error) {
	// if JUnit output is desired, set up temporary file to which raw output is written
	var rawFile *os.File
	rawWriter := ioutil.Discard
	if params.junitOutputPath != "" {
		var err error
		rawFile, err = ioutil.TempFile("", "")
		if err != nil {
			return errors.Wrapf(err, "failed to create temporary file")
		}
		rawWriter = rawFile
		defer func() {
			if err := os.Remove(rawFile.Name()); err != nil && rErr == nil {
				rErr = errors.Wrapf(err, "failed to remove temporary file %s in defer", rawFile.Name())
			}
		}()
	}

	// run the test function
	params.verbose = params.verbose || params.junitOutputPath != ""
	cmder, err := supplier(g.cmd)
	if err != nil {
		return errors.Wrapf(err, "failed to create command %s", g.cmd.Name())
	}
	failedPkgs, err := g.runFunc(cmder, pkgs, params, rawWriter, wd)

	// close raw file
	if rawFile != nil {
		if err := rawFile.Close(); err != nil {
			return errors.Wrapf(err, "failed to close file %s", rawFile.Name())
		}
	}

	if err != nil && err.Error() != "exit status 1" {
		// only re-throw if error is not "exit status 1", since those errors are generally recoverable
		return err
	}

	if params.junitOutputPath != "" {
		// open raw file for reading
		var err error
		rawFile, err = os.Open(rawFile.Name())
		if err != nil {
			return errors.Wrapf(err, "failed to open temporary file %s for reading", rawFile.Name())
		}

		if err := runGoJUnitReport(supplier, wd, rawFile, params.junitOutputPath); err != nil {
			return err
		}
	}

	if len(failedPkgs) > 0 {
		return fmt.Errorf(failedPkgsErrorMsg(failedPkgs))
	}

	return nil
}

type paramCreatorFunc func(ctx cli.Context) testCtxParams

func testParamCreator(ctx cli.Context) testCtxParams {
	return testCtxParams{
		junitOutputPath: cmd.JUnitOutputPath(ctx),
		race:            cmd.Race(ctx),
		stdout:          ctx.App.Stdout,
		tags:            cmd.Tags(ctx),
		verbose:         cmd.Verbose(ctx),
	}
}

func coverageParamCreator(ctx cli.Context) testCtxParams {
	param := testParamCreator(ctx)
	param.coverageOutPath = ctx.String(coverageOutputPathFlag)
	return param
}

type testCtxParams struct {
	stdout          io.Writer
	coverageOutPath string
	junitOutputPath string
	verbose         bool
	race            bool
	tags            []string
}

func longestPkgNameLen(pkgPaths []string, cmdWd string) int {
	longestPkgLen := 0
	for _, currPkgPath := range pkgPaths {
		pkgName, err := pkgpath.NewRelPkgPath(currPkgPath, cmdWd).GoPathSrcRel()
		if err == nil && len(pkgName) > longestPkgLen {
			longestPkgLen = len(pkgName)
		}
	}
	return longestPkgLen
}

type runTestFunc func(cmder amalgomated.Cmder, pkgs []string, params testCtxParams, w io.Writer, wd string) ([]string, error)

func runGoTest(cmder amalgomated.Cmder, pkgs []string, params testCtxParams, w io.Writer, wd string) ([]string, error) {
	// make test output verbose
	if params.verbose {
		cmder = amalgomated.CmderWithPrependedArgs(cmder, "-v")
	}
	if params.race {
		cmder = amalgomated.CmderWithPrependedArgs(cmder, "-race")
	}
	return executeTestCommand(cmder.Cmd(pkgs, wd), params.stdout, w, longestPkgNameLen(pkgs, wd))
}

func runGoTestCover(cmder amalgomated.Cmder, pkgs []string, params testCtxParams, w io.Writer, wd string) (rFailedTests []string, rErr error) {
	// create combined output file
	outputFile, err := os.Create(params.coverageOutPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create specified output file for coverage at %s", params.coverageOutPath)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			fmt.Printf("%+v\n", errors.Wrapf(err, "failed to close file %s in defer", outputFile))
		}
	}()

	// create temporary directory for individual coverage profiles
	tmpDir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create temporary directory for coverage output")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to remove temporary directory %s in defer", tmpDir)
		}
	}()

	isFirstPackage := true
	var failedTests []string
	// currently can only run one package at a time
	for _, currPkg := range pkgs {
		// if error existed, add package to failed tests
		longestPkgNameLen := longestPkgNameLen(pkgs, wd)
		failedPkgs, currPkgCoverageFilePath, err := coverSinglePkg(cmder, params.stdout, w, wd, params.verbose, params.race, currPkg, tmpDir, longestPkgNameLen)
		if err != nil {
			failedTests = append(failedTests, failedPkgs...)
		}

		if err := appendSingleCoverageOutputToCombined(currPkgCoverageFilePath, isFirstPackage, outputFile); err != nil {
			return nil, err
		}
		isFirstPackage = false
	}

	return failedTests, err
}

func appendSingleCoverageOutputToCombined(singleCoverageFilePath string, isFirstPkg bool, combinedOutput io.Writer) (rErr error) {
	singlePkgCoverageFile, err := os.Open(singleCoverageFilePath)
	if err != nil {
		return errors.Wrapf(err, "failed to open file %s", singleCoverageFilePath)
	}

	defer func() {
		if err := singlePkgCoverageFile.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close file %s in defer", singleCoverageFilePath)
		}
	}()

	// append current output to combined output file
	br := bufio.NewReader(singlePkgCoverageFile)
	if !isFirstPkg {
		// if this is not the first package, skip the first line (it contains the coverage mode)
		if _, err := br.ReadString('\n'); err != nil {
			// do nothing
		}
	}
	if _, err := io.Copy(combinedOutput, br); err != nil {
		return errors.Wrapf(err, "failed to write output to writer")
	}

	return nil
}

// coverSinglePkgs runs the cover command on a single package. The raw output of the command written to the provided
// writer. The coverage profile for the file is written to a temporary file within the provided directory. The function
// returns the names of any packages that failed (should be either empty or a slice containing the package name of the
// package that was covered), the location that the coverage profile for this package was written and an error value.
func coverSinglePkg(cmder amalgomated.Cmder, stdout io.Writer, rawWriter io.Writer, cmdWd string, verbose, race bool, currPkg, tmpDir string, longestPkgNameLen int) (rFailedPkgs []string, rTmpFile string, rErr error) {
	currTmpFile, err := ioutil.TempFile(tmpDir, "")
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to create temporary file for coverage for package %s", currPkg)
	}
	defer func() {
		if err := currTmpFile.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close file %s in defer", currTmpFile.Name())
		}
	}()

	// make test output verbose and enable coverage
	var prependedArgs []string
	if verbose {
		prependedArgs = append(prependedArgs, "-v")
	}
	if race {
		prependedArgs = append(prependedArgs, "-race")
	}
	prependedArgs = append(prependedArgs, "-covermode=count", "-coverprofile="+currTmpFile.Name())
	wrappedCmder := amalgomated.CmderWithPrependedArgs(cmder, prependedArgs...)

	// execute test for package
	failedPkgs, err := executeTestCommand(wrappedCmder.Cmd([]string{currPkg}, cmdWd), stdout, rawWriter, longestPkgNameLen)
	return failedPkgs, currTmpFile.Name(), err
}

func runGoJUnitReport(supplier amalgomated.CmderSupplier, cmdWd string, rawOutputReader io.Reader, junitOutputPath string) error {
	cmder, err := supplier(goJUnitReportCmd)
	if err != nil {
		return errors.Wrapf(err, "failed to create runner for gojunitreport")
	}

	var goVersionArgs []string
	if versionOutput, err := exec.Command("go", "version").CombinedOutput(); err == nil {
		if parts := strings.Split(string(versionOutput), " "); len(parts) >= 3 {
			// expect output to be of form "go version go1.8 darwin/amd64", so get element 2
			goVersionArgs = append(goVersionArgs, fmt.Sprintf("--go-version=%s", parts[2]))
		}
	}

	execCmd := cmder.Cmd(goVersionArgs, cmdWd)
	execCmd.Stdin = bufio.NewReader(rawOutputReader)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "failed to run gojunitreport")
	}

	if err := ioutil.WriteFile(junitOutputPath, output, 0644); err != nil {
		return errors.Wrapf(err, "failed to write output to path %s", junitOutputPath)
	}
	return nil
}

// createPlaceholderTestFiles creates placeholder test files in any of the provided packages that do not already contain
// test files and returns a slice that contains the created files. If this function returns an error, it will attempt to
// remove any of the placeholder files that it created before doing so. The generated files will have the name
// "tmp_placeholder_test.go" and will have a package clause that matches the name of the other go files in the
// directory.
func createPlaceholderTestFiles(pkgs []string, wd string) ([]string, error) {
	var placeholderFiles []string

	for _, currPkg := range pkgs {
		currPath := path.Join(wd, currPkg)
		infos, err := ioutil.ReadDir(currPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get directory information for package %s", currPath)
		}

		pkgHasTest := false
		for _, currFileInfo := range infos {
			if !currFileInfo.IsDir() && strings.HasSuffix(currFileInfo.Name(), "_test.go") {
				pkgHasTest = true
				break
			}
		}

		// no test present -- get package name and write temporary placeholder file
		if !pkgHasTest {
			parsedPkgs, err := parser.ParseDir(token.NewFileSet(), currPath, nil, parser.PackageClauseOnly)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse packages from directory %s", currPath)
			}

			if len(parsedPkgs) > 0 {
				// get package name (should only be one since there are no tests in this directory and
				// go requires one package per directory, but will work even if that is not the case)
				pkgName := ""

				// attempt to determine package by importing package using the default build context
				// first. Useful because this will take build constraints into account (for example, if
				// there are Go files in the directory with a "// +build ignore" tag, those should not
				// be considered for generating the package name).
				if importedPkg, err := build.Default.ImportDir(currPath, build.ImportComment); err == nil {
					pkgName = importedPkg.Name
				} else {
					var sortedKeys []string
					for k := range parsedPkgs {
						sortedKeys = append(sortedKeys, k)
					}
					sort.Strings(sortedKeys)
					for _, currName := range sortedKeys {
						pkgName = currName
						break
					}
				}

				currPlaceholderFile := path.Join(currPath, "tmp_placeholder_test.go")

				if err := ioutil.WriteFile(currPlaceholderFile, placeholderTestFileBytes(pkgName), 0644); err != nil {
					// if write fails, clean up files that were already written before returning
					for _, currFileToClean := range placeholderFiles {
						if err := os.Remove(currFileToClean); err != nil {
							fmt.Printf("failed to remove file %s: %v\n", currFileToClean, err)
						}
					}
					return nil, errors.Wrapf(err, "failed to write placeholder file %s", currPlaceholderFile)
				}

				placeholderFiles = append(placeholderFiles, currPlaceholderFile)
			}
		}
	}

	return placeholderFiles, nil
}

const placeholderTemplate = `package {{package}}
// temporary placeholder test file created by gunit
`

func placeholderTestFileBytes(pkgName string) []byte {
	return []byte(strings.Replace(placeholderTemplate, "{{package}}", pkgName, -1))
}

// executeTestCommand executes the provided command. The output produced by the command's Stdout and Stderr calls are
// processed as they are written and an aligned version of the output is written to the Stdout of the current process.
// The "longestPkgNameLen" parameter specifies the longest package name (used to align the console output). This
// function returns a slice that contains the packages that had test failures (output line started with "FAIL"). The
// error value will contain any error that was encountered while executing the command, including if the command
// executed successfully but any tests failed. In either case, the packages that encountered errors will also be
// returned.
func executeTestCommand(execCmd *exec.Cmd, stdout io.Writer, rawOutputWriter io.Writer, longestPkgNameLen int) (rFailedPkgs []string, rErr error) {
	bw := bufio.NewWriter(rawOutputWriter)

	// stream output to Stdout
	multiWriter := multiWriter{
		consoleWriter:     stdout,
		rawOutputWriter:   bw,
		failedPkgs:        []string{},
		longestPkgNameLen: longestPkgNameLen,
	}

	// flush buffered writer at the end of the function
	defer func() {
		if err := bw.Flush(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to flush buffered writer in defer")
		}
	}()

	// set Stdout and Stderr of command to multiwriter
	execCmd.Stdout = &multiWriter
	execCmd.Stderr = &multiWriter

	// run command (which will print its Stdout and Stderr to the Stdout of current process) and return output
	err := execCmd.Run()
	return multiWriter.failedPkgs, err
}

type multiWriter struct {
	consoleWriter     io.Writer
	rawOutputWriter   io.Writer
	failedPkgs        []string
	longestPkgNameLen int
}

var setupFailedRegexp = regexp.MustCompile(`(^FAIL\t.+) (\[setup failed\]$)`)

func (w *multiWriter) Write(p []byte) (int, error) {
	// write unaltered output to file
	n, err := w.rawOutputWriter.Write(p)
	if err != nil {
		return n, err
	}

	lines := strings.Split(string(p), "\n")
	for i, currLine := range lines {
		// test output for valid case always starts with "Ok" or "FAIL"
		if strings.HasPrefix(currLine, "ok") || strings.HasPrefix(currLine, "FAIL") {
			if setupFailedRegexp.MatchString(currLine) {
				// if line matches "setup failed" output, modify output to conform to expected style
				// (namely, replace space between package name and "[setup failed]" with a tab)
				currLine = setupFailedRegexp.ReplaceAllString(currLine, "$1\t$2")
			}

			// split into at most 4 parts
			fields := strings.SplitN(currLine, "\t", 4)

			// valid test lines have at least 3 parts: "[ok|FAIL]\t[pkgName]\t[time|no test files]"
			if len(fields) >= 3 {
				currPkgName := strings.TrimSpace(fields[1])
				lines[i] = alignLine(fields, w.longestPkgNameLen)
				// append package name to failures list if this was a failure
				if strings.HasPrefix(currLine, "FAIL") {
					w.failedPkgs = append(w.failedPkgs, currPkgName)
				}
			}
		}
	}

	// write formatted version to console writer
	if n, err := w.consoleWriter.Write([]byte(strings.Join(lines, "\n"))); err != nil {
		return n, err
	}

	// n and err are from the unaltered write to the rawOutputWriter
	return n, err
}

// alignLine returns a string where the length of the second field (fields[1]) is padded with spaces to make its length
// equal to the value of maxPkgLen and the fields are joined with tab characters. Assuming that the first field is
// always the same length, this method ensures that the third field will always be aligned together for any fixed value
// of maxPkgLen.
func alignLine(fields []string, maxPkgLen int) string {
	currPkgName := fields[1]
	repeat := maxPkgLen - len(currPkgName)
	if repeat < 0 {
		// this should not occur under normal circumstances. However, it appears that it is possible if tests
		// create test packages in the directory structure while tests are already running. If such a case is
		// encountered, having output that isn't aligned optimally is better than crashing, so set repeat to 0.
		repeat = 0
	}
	fields[1] = currPkgName + strings.Repeat(" ", repeat)
	return strings.Join(fields, "\t")
}

func failedPkgsErrorMsg(failedPkgs []string) string {
	numFailedPkgs := len(failedPkgs)
	outputParts := append([]string{fmt.Sprintf("%d %v had failing tests:", numFailedPkgs, plural(numFailedPkgs, "package", "packages"))}, failedPkgs...)
	return strings.Join(outputParts, "\n\t")
}

func plural(num int, singular, plural string) string {
	if num == 1 {
		return singular
	}
	return plural
}
