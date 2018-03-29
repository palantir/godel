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

package generator

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/palantir/pkg/signals"
	"github.com/pkg/errors"
)

type Params struct {
	TagPrefix            string
	RunDockerBuild       bool
	SuppressDockerOutput bool
	StartStep            int
	EndStep              int
	LeaveGeneratedFiles  bool
}

func Generate(inputDir, outputDir, baseImage string, params Params, stdout io.Writer) error {
	inputFiles, err := getInputFilesFromDir(inputDir)
	if err != nil {
		return err
	}

	fmt.Fprintf(stdout, "Found %d template file(s)\n", len(inputFiles))
	startIdx := 0
	if params.StartStep != -1 {
		startIdx = findIdxWithOrdering(params.StartStep, inputFiles)
		if startIdx == -1 {
			return errors.Errorf("could not find specified start step %d in %v", params.StartStep, inputFiles)
		}
	}
	endIdx := len(inputFiles) - 1
	if params.EndStep != -1 {
		endIdx = findIdxWithOrdering(params.EndStep, inputFiles)
		if endIdx == -1 {
			return errors.Errorf("could not find specified end step %d in %v", params.StartStep, inputFiles)
		}
	}

	numFiles := endIdx - startIdx + 1
	currCount := 1
	fmt.Fprintf(stdout, "Processing %d template file(s) starting at number %d and ending at number %d\n", numFiles, inputFiles[startIdx].Ordering, inputFiles[endIdx].Ordering)

	for idx, inputFile := range inputFiles {
		if idx < startIdx || idx > endIdx {
			continue
		}

		err := func() error {
			fmt.Fprintf(stdout, "Processing %s (%d/%d)\n", inputFile.TemplateFileName, currCount, numFiles)

			var fromImage string
			if idx == 0 {
				fromImage = baseImage
			} else {
				fromImage = inputFiles[idx-1].DockerTag(params.TagPrefix)
			}
			fileWithContent, err := readTemplateFile(inputDir, inputFile)
			if err != nil {
				return errors.Wrapf(err, "failed to read template file")
			}

			fmt.Fprintln(stdout, "Writing output files...")
			currOutputDir, err := writeOutputFiles(outputDir, fileWithContent, fromImage)
			if err != nil {
				return errors.Wrapf(err, "failed to write output files for %s", inputFile.TemplateFileName)
			}
			// run the rest of the logic in a wrapped function to allow deferral of removing temporary directory
			if err := func() (rErr error) {
				if !params.LeaveGeneratedFiles {
					cleanupCtx, cancel := signals.ContextWithShutdown(context.Background())
					cleanupDone := make(chan struct{})
					defer func() {
						cancel()
						<-cleanupDone
					}()
					go func() {
						select {
						case <-cleanupCtx.Done():
							if err := os.RemoveAll(currOutputDir); err != nil && rErr == nil {
								rErr = errors.Wrapf(err, "failed to remove output directory")
							}
						}
						cleanupDone <- struct{}{}
					}()
				}
				if !params.RunDockerBuild {
					return nil
				}
				fmt.Fprintln(stdout, "Running Docker build...")
				dockerBuildOutput, err := runDockerBuild(currOutputDir, inputFile.DockerTag(params.TagPrefix), params.SuppressDockerOutput, stdout)
				if err != nil {
					return errors.Wrapf(err, "docker build failed for %s", inputFile.TemplateFileName)
				}
				bashRunCmds := parseBashRunCmdFromOutput(dockerBuildOutput)

				if len(bashRunCmds) != len(fileWithContent.ParsedContent.TutorialCodeParts) {
					return errors.Errorf("number of command outputs did not match number of tutorial code parts: %d != %d", len(bashRunCmds), len(fileWithContent.ParsedContent.TutorialCodeParts))
				}

				var codeToRenderToTemplate []string
				for i := range fileWithContent.ParsedContent.TutorialCodeParts {
					// work-around: replace command to execute parsed from output with the original command. This will
					// ensure that any modifications to the command made by the flags/options are not included in the
					// template output.
					bashRunCmds[i].cmd = fileWithContent.ParsedContent.TutorialCodeParts[i].Code
					codeToRenderToTemplate = append(codeToRenderToTemplate, bashRunCmds[i].String())
				}
				renderedBytes, err := fileWithContent.ParsedContent.Render(codeToRenderToTemplate)
				if err != nil {
					return errors.Wrapf(err, "failed to render parsed content for %s", inputFile.TemplateFileName)
				}
				renderedFilePath := path.Join(outputDir, inputFile.OutputRenderedFileName())
				if err := ioutil.WriteFile(renderedFilePath, renderedBytes, 0644); err != nil {
					return errors.Wrapf(err, "failed to write rendered content for %s", inputFile.TemplateFileName)
				}
				return nil
			}(); err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return errors.Wrapf(err, "failed running task for template %s (%d/%d)", inputFile.TemplateFileName, currCount, numFiles)
		}
		currCount++
	}
	return nil
}

func runDockerBuild(workDir, tag string, suppressDockerOutput bool, stdout io.Writer) (string, error) {
	cmd := exec.Command("docker", "build", "--no-cache", "-t", tag, ".")
	cmd.Dir = workDir

	outputBuf := &bytes.Buffer{}
	writers := []io.Writer{
		outputBuf,
	}
	if !suppressDockerOutput {
		writers = append(writers, stdout)
	}
	mw := io.MultiWriter(writers...)
	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Run(); err != nil {
		errMsg := fmt.Sprintf("command %v failed", cmd.Args)
		if suppressDockerOutput {
			errMsg += " with output " + outputBuf.String()
		}
		return "", errors.Wrapf(err, errMsg)
	}
	return outputBuf.String(), nil
}

func findIdxWithOrdering(wantOrdering int, files []inputFile) int {
	idx := -1
	for i, file := range files {
		if file.Ordering == wantOrdering {
			idx = i
			break
		}
	}
	return idx
}

func newInputFile(fileName string) (inputFile, error) {
	submatch := templateFileRegexp.FindStringSubmatch(fileName)
	if submatch == nil {
		return inputFile{}, errors.Errorf("input %q does not match required format", fileName)
	}
	ordering, err := strconv.Atoi(submatch[1])
	if err != nil {
		return inputFile{}, errors.Wrapf(err, "failed to parse ordering component of %q", fileName)
	}
	name := submatch[2]
	var originalExtension string
	if lastDot := strings.LastIndex(submatch[2], "."); lastDot != -1 {
		name = submatch[2][:lastDot]
		originalExtension = submatch[2][lastDot:]
	}
	return inputFile{
		TemplateFileName:  fileName,
		Ordering:          ordering,
		Name:              name,
		OriginalExtension: originalExtension,
	}, nil
}

type inputFile struct {
	TemplateFileName  string
	Ordering          int
	Name              string
	OriginalExtension string
}

func (f *inputFile) OutputDirName() string {
	return fmt.Sprintf("%d_%s", f.Ordering, f.Name)
}

func (f *inputFile) OutputScriptFileName() string {
	return fmt.Sprintf("run-%s.sh", f.Name)
}

func (f *inputFile) OutputRenderedFileName() string {
	return f.Name + f.OriginalExtension
}

func (f *inputFile) DockerTag(tagPrefix string) string {
	return fmt.Sprintf("%s:%s", tagPrefix, f.Name)
}

var templateFileRegexp = regexp.MustCompile(`^(\d+)_(.+)` + regexp.QuoteMeta(".tmpl") + "$")

func getInputFilesFromDir(inputDir string) ([]inputFile, error) {
	fis, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read directory")
	}
	var fileNames []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		fileNames = append(fileNames, fi.Name())
	}
	return getInputFiles(fileNames)
}

func getInputFiles(fileNames []string) ([]inputFile, error) {
	var inputFiles []inputFile

	// values are the template names
	orderings := make(map[int][]string)
	names := make(map[string][]string)

	for _, fileName := range fileNames {
		currInputFile, err := newInputFile(fileName)
		if err != nil {
			continue
		}
		inputFiles = append(inputFiles, currInputFile)
		orderings[currInputFile.Ordering] = append(orderings[currInputFile.Ordering], currInputFile.TemplateFileName)
		names[currInputFile.Name] = append(names[currInputFile.Name], currInputFile.TemplateFileName)
	}
	// sort by ordering
	sort.Slice(inputFiles, func(i, j int) bool {
		return inputFiles[i].Ordering < inputFiles[j].Ordering
	})

	// verify that ordering values are unique
	var sortedOrderings []int
	for k := range orderings {
		sortedOrderings = append(sortedOrderings, k)
	}
	sort.Ints(sortedOrderings)
	for _, k := range sortedOrderings {
		if len(orderings[k]) > 1 {
			return nil, errors.Errorf("multiple inputs have the ordering value %d: %v", k, orderings[k])
		}
	}

	// verify that name values are unique
	var sortedNames []string
	for k := range names {
		sortedNames = append(sortedNames, k)
	}
	sort.Strings(sortedNames)
	for _, k := range sortedNames {
		if len(names[k]) > 1 {
			return nil, errors.Errorf("multiple inputs have the name %q: %v", k, names[k])
		}
	}

	return inputFiles, nil
}
