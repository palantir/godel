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

package cmd

import (
	"bytes"
	"io"
	"path"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/gonform/config"
	"github.com/palantir/godel/apps/gonform/params"
)

var (
	goFmt     = Library.MustNewCmd("gofmt")
	ptImports = Library.MustNewCmd("ptimports")
)

func RunAllCommand(supplier amalgomated.CmderSupplier) cli.Command {
	return cli.Command{
		Name:  "runAll",
		Usage: "Run all format commands on Go files",
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return DoRunAll(ctx.Slice(filesParamName), ctx, supplier, wd)
		},
		Flags: []flag.Flag{
			FilesParam,
		},
	}
}

func GoFmtCommand(supplier amalgomated.CmderSupplier) cli.Command {
	return formatCommand(goFmt, "Run ptimports on Go files", supplier)
}

func PTImportsCommand(supplier amalgomated.CmderSupplier) cli.Command {
	return formatCommand(ptImports, "Run ptimports on Go files", supplier)
}

func DoRunAll(filesParam []string, ctx cli.Context, supplier amalgomated.CmderSupplier, wd string) error {
	cmds := []amalgomated.Cmd{
		goFmt,
		ptImports,
	}
	params, err := createParams(filesParam, ctx, wd)
	if err != nil {
		return err
	}
	for _, currCmd := range cmds {
		err := doFormat(currCmd, params, ctx, supplier, wd)
		if err != nil {
			return err
		}
	}
	return nil
}

func formatCommand(cmd amalgomated.Cmd, usage string, supplier amalgomated.CmderSupplier) cli.Command {
	return cli.Command{
		Name:  cmd.Name(),
		Usage: usage,
		Action: func(ctx cli.Context) error {
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}
			return doSingleFormat(cmd, ctx.Slice(filesParamName), ctx, supplier, wd)
		},
		Flags: []flag.Flag{
			FilesParam,
		},
	}
}

type formatterParams struct {
	Formatters map[string]params.Formatter
	Files      []string
	List       bool
	Verbose    bool
}

func doSingleFormat(cmd amalgomated.Cmd, filesParam []string, ctx cli.Context, supplier amalgomated.CmderSupplier, wd string) error {
	params, err := createParams(filesParam, ctx, wd)
	if err != nil {
		return err
	}
	return doFormat(cmd, params, ctx, supplier, wd)
}

func createParams(filesParam []string, ctx cli.Context, wd string) (formatterParams, error) {
	cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
	if err != nil {
		return formatterParams{}, errors.Wrapf(err, "failed to load configuration")
	}

	files, err := getFiles(filesParam, cfg.Exclude, wd)
	if err != nil {
		return formatterParams{}, errors.Wrapf(err, "failed to get files")
	}

	return formatterParams{
		Formatters: cfg.Formatters,
		Files:      files,
		List:       ctx.Bool(listFlagName),
		Verbose:    ctx.Bool(verboseFlagName),
	}, nil
}

func doFormat(cmd amalgomated.Cmd, params formatterParams, ctx cli.Context, supplier amalgomated.CmderSupplier, wd string) error {
	if len(params.Files) == 0 {
		return nil
	}

	supplier = amalgomated.SupplierWithPrependedArgs(supplier, func(cmd amalgomated.Cmd) []string {
		formatArg := "-w"
		if params.List {
			formatArg = "-l"
		}
		return append(params.Formatters[cmd.Name()].Args, formatArg)
	})

	combinedBuf := bytes.Buffer{}
	stdoutMultiWriter := io.MultiWriter(ctx.App.Stdout, &combinedBuf)
	stderrMultiWriter := io.MultiWriter(ctx.App.Stderr, &combinedBuf)

	cmder, err := supplier(cmd)
	if err != nil {
		return errors.Wrapf(err, "failed to create command %s", cmd.Name())
	}

	if params.Verbose {
		ctx.Printf("Running %s...\n", cmd.Name())
	}
	execCmd := cmder.Cmd(params.Files, wd)
	execCmd.Stdout = stdoutMultiWriter
	execCmd.Stderr = stderrMultiWriter

	if err := execCmd.Run(); err != nil {
		return errors.Wrapf(err, "command %s failed", cmd.Name())
	} else if params.List && len(combinedBuf.Bytes()) != 0 {
		if params.Verbose {
			return errors.Errorf("%s failed: unformatted Go files found", cmd.Name())
		}
		return errors.Errorf("")
	}
	return nil
}

func getFiles(filesParam []string, exclude matcher.Matcher, wd string) ([]string, error) {
	var files []string
	var err error
	if len(filesParam) == 0 {
		// exclude entries specified by the configuration
		files, err = matcher.ListFiles(wd, matcher.Name(`.*\.go`), exclude)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list *.go files in %d", wd)
		}
	} else {
		// filter arguments based on exclude config
		files = make([]string, 0, len(filesParam))
		for _, currPath := range filesParam {
			if !exclude.Match(currPath) {
				files = append(files, currPath)
			}

		}
	}
	files = toAbsPaths(wd, files)

	if len(files) == 0 {
		// no files to process after exclusions applied
		return nil, nil
	}

	return files, nil
}

func toAbsPaths(basePath string, relPaths []string) []string {
	paths := make([]string, len(relPaths))
	for i, currPath := range relPaths {
		paths[i] = path.Join(basePath, currPath)
	}
	return paths
}
