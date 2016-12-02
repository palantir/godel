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

package matcher

// NamesPathsCfg is a configuration object that defines a list of names and paths that should be used to construct a
// Matcher. The returned Matcher will match any name or path specified in the configuration.
type NamesPathsCfg struct {
	Names []string `yaml:"names" json:"names"`
	Paths []string `yaml:"paths" json:"paths"`
}

// Add appends the names and paths specified in the provided NamesPathsCfg to those in the receiver.
func (c *NamesPathsCfg) Add(cfg NamesPathsCfg) {
	c.Names = append(c.Names, cfg.Names...)
	c.Paths = append(c.Paths, cfg.Paths...)
}

// Empty returns true if the configuration is empty. If this function returns true, it indicates that the semantic
// meaning of the configuration is the same as it not being provided/specified at all.
func (c *NamesPathsCfg) Empty() bool {
	return len(c.Names) == 0 && len(c.Paths) == 0
}

// Matcher returns a Matcher constructed from the configuration. The Matcher returns true if it matches any of the
// names or paths in the configuration.
func (c *NamesPathsCfg) Matcher() Matcher {
	return Any(Name(c.Names...), Path(c.Paths...))
}
