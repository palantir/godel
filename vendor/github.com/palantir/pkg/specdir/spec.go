// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package specdir

import (
	"fmt"
	"os"
	"path"
)

type PathType bool

const (
	DirPath  = PathType(false)
	FilePath = PathType(true)
)

func NewLayoutSpec(root FileNodeProvider, rootPartOfSpec bool) LayoutSpec {
	return &layoutSpec{
		root:           root.fileNode(),
		rootPartOfSpec: rootPartOfSpec,
	}
}

// LayoutSpec represents the specification of a layout.
type LayoutSpec interface {
	FileNodeProvider

	CreateDirectoryStructure(root string, values TemplateValues, includeOptional bool) error
	RootDirName(values TemplateValues) string
	Paths(values TemplateValues, includeOptional bool) []string
	// Validate validates the provided root directory against this specification using the provided template values. If the
	// provided directory matches the specification, the function returns nil. If the function encounters an error while
	// trying to perform a verification or if the provided structure does not match the specification, an error is returned.
	Validate(root string, values TemplateValues) error

	rootNode() *fileNode
	rootIsPartOfSpec() bool
	getAliasValues(templateValues map[string]string) map[string]string
}

type layoutSpec struct {
	// must be a non-optional directory
	root *fileNode
	// rootNotPartOfSpec is true only if the root directory itself is not part of the specification. If this is
	// true, then the specification only enforces that the content of the root node match the content of the
	// provided root (the roots themselves can have different names).
	rootPartOfSpec bool
}

// CreateDirectoryStructure creates the directory structure necessary to make the provided root directory match the
// specification. This method only creates directories (it will not attempt to create any files). This method will not
// overwrite any existing files or directories, and if an existing file conflicts with a directory required by this
// spec, the operation will fail. If the "includeOptional" parameter is true, directories that are marked as "optional"
// will be created; otherwise, optional directories will be omitted. Returns an error if the operation does not succeed.
// If an error is returned, the method will make a best effort to remove the directories that it created. If the
// specification dictates that the root directory is part of the spec and the provided root directory does not match the
// specification, an error is returned.
func (s *layoutSpec) CreateDirectoryStructure(root string, values TemplateValues, includeOptional bool) error {
	var missingKeys []string
	for _, currKey := range s.getNameTemplateKeys() {
		if _, ok := values[currKey]; !ok {
			missingKeys = append(missingKeys, currKey)
		}
	}
	if len(missingKeys) > 0 {
		return fmt.Errorf("required template keys were missing: got %v, missing %v", values, missingKeys)
	}

	if s.rootPartOfSpec {
		// verify provided path matches root directory
		if err := verifyPath("", root, s.root.name.name(values), DirPath, false); err != nil {
			return err
		}
	}

	for _, c := range s.root.children {
		if err := c.fileNode().createDirectoryStructure(root, values, includeOptional); err != nil {
			return err
		}
	}

	return nil
}

func (s *layoutSpec) rootNode() *fileNode {
	return s.root
}

func (s *layoutSpec) rootIsPartOfSpec() bool {
	return s.rootPartOfSpec
}

func (s *layoutSpec) RootDirName(values TemplateValues) string {
	return s.root.name.name(values)
}

func (s *layoutSpec) Paths(values TemplateValues, includeOptional bool) []string {
	var paths []string

	currPath := ""
	if s.rootPartOfSpec {
		currPath = s.root.name.name(values)
		paths = append(paths, currPath)
	}

	// add all children
	paths = append(paths, s.root.pathsForChildrenOfDir(currPath, values, includeOptional)...)

	return paths
}

func (s *layoutSpec) fileNode() *fileNode {
	return s.root
}

func (s *layoutSpec) getAliasValues(templateValues map[string]string) map[string]string {
	aliasValues := make(map[string]string)

	var currPath []*fileNode
	if s.rootPartOfSpec {
		currPath = []*fileNode{s.root}

		if s.root.alias != "" {
			aliasValues[s.root.alias] = fileNodePath(currPath).getPath(templateValues)
		}
	}

	s.root.addAliasValuesForDirNodeChildren(aliasValues, currPath, templateValues)

	return aliasValues
}

// Validate validates the provided root directory against this specification using the provided template values. If the
// provided directory matches the specification, the function returns nil. If the function encounters an error while
// trying to perform a verification or if the provided structure does not match the specification, an error is returned.
func (s *layoutSpec) Validate(root string, values TemplateValues) error {
	if s.rootPartOfSpec {
		// verify provided path matches root directory
		if err := verifyPath("", root, s.root.name.name(values), DirPath, false); err != nil {
			return err
		}
	}

	// verify all child nodes match
	return s.root.verifyLayoutForDir(root, "", values)
}

func (s *layoutSpec) getNameTemplateKeys() []string {
	return s.root.getNameTemplateKeys()
}

func getTemplateKeysFromName(n NodeName) []string {
	var names []string
	switch t := n.(type) {
	case templateName:
		names = append(names, string(t))
	case compositeName:
		for _, currName := range t {
			names = append(names, getTemplateKeysFromName(currName)...)
		}
	}
	return names
}

func verifyPath(rootDirPath, pathFromRootDir, expectedName string, pathType PathType, optional bool) error {
	if path.Base(pathFromRootDir) != expectedName {
		return fmt.Errorf("%s is not a path to %s", pathFromRootDir, expectedName)
	}

	pathInfo, err := os.Stat(path.Join(rootDirPath, pathFromRootDir))
	if err != nil {
		if os.IsNotExist(err) {
			if !optional {
				return fmt.Errorf("%s does not exist", path.Join(path.Base(rootDirPath), pathFromRootDir))
			}
			// path does not exist, but it is optional so is okay
			return nil
		}
		return fmt.Errorf("failed to stat %s", path.Join(path.Base(rootDirPath), pathFromRootDir))
	} else if currIsDir := pathInfo.IsDir(); currIsDir == bool(pathType) {
		return fmt.Errorf("isDir for %s returned wrong value: expected %v, was %v", path.Join(path.Base(rootDirPath), pathFromRootDir), !pathType, currIsDir)
	}

	return nil
}
