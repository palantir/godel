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

// NodeName represents the name of a node and returns a value for its name given a TemplateValues object. An empty
// string is returned if a name cannot be generated based on the input.
type NodeName interface {
	name(values TemplateValues) string
}

type fileNode struct {
	name     NodeName
	pathType PathType
	children []FileNodeProvider
	optional bool
	alias    string
}

type templateName string

func (n templateName) name(values TemplateValues) string {
	return values[string(n)]
}

func TemplateName(key string) NodeName {
	return templateName(key)
}

type literalName string

func LiteralName(name string) NodeName {
	return literalName(name)
}

func (n literalName) name(values TemplateValues) string {
	return string(n)
}

type compositeName []NodeName

func (n compositeName) name(values TemplateValues) string {
	combinedName := ""
	for _, currName := range n {
		combinedName += currName.name(values)
	}
	return combinedName
}

func CompositeName(names ...NodeName) NodeName {
	return compositeName(names)
}

type TemplateValues map[string]string

type FileNodeProvider interface {
	fileNode() *fileNode
}

func (n *fileNode) fileNode() *fileNode {
	return n
}

func (n *fileNode) createDirectoryStructure(parentDir string, values TemplateValues, includeOptional bool) error {
	if n.pathType == DirPath && (!n.optional || includeOptional) {
		currPath := path.Join(parentDir, n.name.name(values))
		if err := os.MkdirAll(currPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", currPath, err)
		}

		for _, c := range n.children {
			if err := c.fileNode().createDirectoryStructure(currPath, values, includeOptional); err != nil {
				return err
			}
		}
	}

	return nil
}

type fileNodePath []*fileNode

func (p fileNodePath) getPath(templateValues map[string]string) string {
	var parts []string
	for _, node := range p {
		parts = append(parts, node.name.name(templateValues))
	}
	return path.Join(parts...)
}

// getTemplateNames returns the names of all of the templates in the given path
func (p fileNodePath) getTemplateNames() []string {
	var names []string
	for _, currNode := range p {
		names = append(names, getTemplateKeysFromName(currNode.name)...)
	}
	return names
}

func (n *fileNode) paths(parentPath string, templateValues map[string]string, includeOptional bool) []string {
	var paths []string

	// only add node and children if it is not optional or if optional paths are being included
	if !n.optional || includeOptional {
		currPath := path.Join(parentPath, n.name.name(templateValues))

		// add current node
		paths = append(paths, currPath)

		// add all child nodes
		paths = append(paths, n.pathsForChildrenOfDir(currPath, templateValues, includeOptional)...)
	}

	return paths
}

func (n *fileNode) pathsForChildrenOfDir(currPath string, templateValues map[string]string, includeOptional bool) []string {
	var paths []string
	// add all child nodes
	if n.pathType == DirPath {
		for _, c := range n.children {
			paths = append(paths, c.fileNode().paths(currPath, templateValues, includeOptional)...)
		}
	}
	return paths
}

func (n *fileNode) getAliasValues(pathSoFar fileNodePath, templateValues map[string]string) map[string]string {
	aliasValues := make(map[string]string)

	currPath := append(append([]*fileNode{}, []*fileNode(pathSoFar)...), n)

	if n.alias != "" {
		aliasValues[n.alias] = fileNodePath(currPath).getPath(templateValues)
	}

	n.addAliasValuesForDirNodeChildren(aliasValues, currPath, templateValues)

	return aliasValues
}

func (n *fileNode) addAliasValuesForDirNodeChildren(aliasValues map[string]string, currPath fileNodePath, templateValues map[string]string) {
	if n.pathType == DirPath {
		for _, c := range n.children {
			for key, value := range c.fileNode().getAliasValues(currPath, templateValues) {
				aliasValues[key] = value
			}
		}
	}
}

func (n *fileNode) validate(rootDir, pathFromRoot string, values TemplateValues) error {
	// verify current path
	if err := verifyPath(rootDir, pathFromRoot, n.name.name(values), n.pathType, n.optional); err != nil {
		return err
	}

	// verify all child nodes match
	return n.verifyLayoutForDir(rootDir, pathFromRoot, values)
}

func (n *fileNode) verifyLayoutForDir(rootDir, pathFromRoot string, values TemplateValues) error {
	if n.pathType == DirPath {
		for _, c := range n.children {
			currPath := path.Join(pathFromRoot, c.fileNode().name.name(values))
			if err := c.fileNode().validate(rootDir, currPath, values); err != nil {
				return err
			}
		}
	}
	return nil
}

func (n *fileNode) getNameTemplateKeys() []string {
	namesMap := make(map[string]bool)

	for _, currName := range getTemplateKeysFromName(n.name) {
		namesMap[currName] = true
	}

	for _, c := range n.children {
		for _, currName := range c.fileNode().getNameTemplateKeys() {
			namesMap[currName] = true
		}
	}

	names := make([]string, 0, len(namesMap))
	for currName := range namesMap {
		names = append(names, currName)
	}
	sort.Strings(names)

	return names
}
