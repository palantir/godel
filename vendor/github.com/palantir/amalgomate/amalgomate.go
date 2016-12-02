// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/pkg/errors"
)

func main() {
	os.Exit(newApp().Run(os.Args))
}

const (
	configFlag    = "config"
	outputDirFlag = "output-dir"
	pkgFlag       = "pkg"
)

func newApp() *cli.App {
	app := cli.NewApp(cli.DebugHandler(errorstringer.SingleStack))
	app.Usage = "Re-package main packages into a library package"
	app.Flags = append(app.Flags,
		flag.StringFlag{
			Name:     configFlag,
			Usage:    "configuration file that specifies packages to be amalgomated",
			Required: true,
		},
		flag.StringFlag{
			Name:     outputDirFlag,
			Usage:    "directory in which amalgomated output is written",
			Required: true,
		},
		flag.StringFlag{
			Name:     pkgFlag,
			Usage:    "package name of the amalgomated source that is generated",
			Required: true,
		},
	)
	app.Action = func(ctx cli.Context) error {
		return doRepackage(ctx.String(configFlag), ctx.String(outputDirFlag), ctx.String(pkgFlag))
	}
	return app
}

func doRepackage(cfgPath, outputDir, pkg string) error {
	// read configuration
	config, err := LoadConfig(cfgPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to read configuration from file: %s", cfgPath)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return errors.Wrapf(err, "Failed to ensure that output directory exists: %s", outputDir)
	}

	if !filepath.IsAbs(outputDir) {
		if outputDir, err = filepath.Abs(outputDir); err != nil {
			return errors.Wrapf(err, "Failed to convert output directory to an absolute path: %s", outputDir)
		}
	}

	// repackage main files specified in configuration
	if err := repackage(*config, outputDir); err != nil {
		return errors.Wrapf(err, "Failed to repackage files specified in config file %s", cfgPath)
	}

	amalgomatedOutputDir := path.Join(outputDir, internalDir)

	// write output file that imports and uses repackaged files
	if err := writeOutputGoFile(*config, outputDir, amalgomatedOutputDir, pkg); err != nil {
		return errors.Wrapf(err, "Failed to write output file")
	}

	return nil
}

func writeOutputGoFile(config Config, outputDir, amalgomatedOutputDir, packageName string) error {
	fileSet := token.NewFileSet()

	var template string
	if packageName == "main" {
		template = mainTemplate
	} else {
		template = libraryTemplate
	}

	file, err := parser.ParseFile(fileSet, "", template, parser.ParseComments)
	if err != nil {
		return errors.Wrapf(err, "failed to parse template: %s", template)
	}
	file.Name = ast.NewIdent(packageName)

	if err := addImports(file, fileSet, amalgomatedOutputDir, config); err != nil {
		return errors.Wrap(err, "failed to add imports")
	}
	sortImports(file)

	if err := setVarCompositeLiteralElements(file, "programs", createMapLiteralEntries(config.Pkgs)); err != nil {
		return errors.Wrap(err, "failed to add const elements")
	}

	// write output to in-memory buffer and add import spaces
	var byteBuffer bytes.Buffer
	if err := printer.Fprint(&byteBuffer, fileSet, file); err != nil {
		return errors.Wrap(err, "failed to write output file to buffer")
	}
	outputWithSpaces := addImportSpaces(&byteBuffer, importBreakPaths(file))

	// write output to file
	outputFilePath := path.Join(outputDir, packageName+".go")
	if err := ioutil.WriteFile(outputFilePath, outputWithSpaces, 0644); err != nil {
		return errors.Wrapf(err, "failed to write output to path: %s", outputFilePath)
	}

	return nil
}
