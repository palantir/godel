// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package specdir

import (
	"fmt"
	"os"
	"path"
	"sort"
)

// SpecDir is a concrete application of a LayoutSpec. It is rooted at a directory and has TemplateValues specified.
type SpecDir interface {
	Root() string
	NamedPaths() []string
	Path(name string) string
}

type specDirStruct struct {
	spec        LayoutSpec
	rootDir     string
	values      TemplateValues
	aliasValues map[string]string
}

func (s *specDirStruct) Root() string {
	return s.rootDir
}

func (s *specDirStruct) NamedPaths() []string {
	names := make([]string, 0, len(s.aliasValues))
	for currName := range s.aliasValues {
		names = append(names, currName)
	}
	sort.Strings(names)
	return names
}

func (s *specDirStruct) Path(name string) string {
	if value, ok := s.aliasValues[name]; ok {
		pathRoot := s.rootDir
		if s.spec.rootIsPartOfSpec() {
			pathRoot = path.Dir(s.rootDir)
		}
		return path.Join(pathRoot, value)
	}
	return ""
}

type Mode int

const (
	// SpecOnly specifies that a specification should be created without performing validation or creating paths.
	SpecOnly Mode = iota
	// Validate specifies that an existing directory should be validated using the specification.
	Validate
	// Create specifies that the specification should create a directory structure that matches the specification.
	Create
)

func New(rootDir string, spec LayoutSpec, values TemplateValues, mode Mode) (SpecDir, error) {
	// verify that all required template values were provided
	requiredNames := spec.rootNode().getNameTemplateKeys()
	for _, name := range requiredNames {
		if _, ok := values[name]; !ok {
			return nil, fmt.Errorf("required template %q was missing.\nRequired: %v\nProvided: %v", name, requiredNames, values)
		}
	}

	switch mode {
	case Validate:
		if err := spec.Validate(rootDir, values); err != nil {
			return nil, err
		}
	case Create:
		// if the root directory is part of the specification, verify that the name matches the required name
		if spec.rootIsPartOfSpec() {
			expectedRootDirName := spec.rootNode().name.name(values)
			if path.Base(rootDir) != expectedRootDirName {
				return nil, fmt.Errorf("root directory name %s does not match name required by specification: %v", rootDir, expectedRootDirName)
			}
		}

		// create the root directory if it does not already exist
		if _, err := os.Stat(rootDir); os.IsNotExist(err) {
			if err := os.Mkdir(rootDir, 0755); err != nil {
				return nil, fmt.Errorf("failed to create directory %s: %v", rootDir, err)
			}
		}

		// create the rest of the directory structure
		if err := spec.CreateDirectoryStructure(rootDir, values, false); err != nil {
			return nil, fmt.Errorf("failed to create directory structure: %v", err)
		}
	case SpecOnly:
		// do nothing
	default:
		return nil, fmt.Errorf("unrecognized mode: %v", mode)
	}

	return &specDirStruct{
		spec:        spec,
		rootDir:     rootDir,
		values:      values,
		aliasValues: spec.getAliasValues(values),
	}, nil
}

func Dir(name NodeName, alias string, children ...FileNodeProvider) FileNodeProvider {
	return &fileNode{
		name:     name,
		children: children,
		alias:    alias,
	}
}

func OptionalDir(name NodeName, children ...FileNodeProvider) FileNodeProvider {
	return &fileNode{
		name:     name,
		children: children,
		optional: true,
	}
}

func File(name NodeName, alias string) FileNodeProvider {
	return &fileNode{
		name:     name,
		pathType: FilePath,
		alias:    alias,
	}
}
