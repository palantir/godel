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

package slsspec_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/specdir"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/pkg/slsspec"
)

func TestSLSSpecCreateDirectoryStructureFailMissingValues(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	spec := slsspec.New()
	err = spec.CreateDirectoryStructure(tmp, nil, false)
	require.Error(t, err)

	assert.Regexp(t, regexp.MustCompile("required template keys were missing: got map[], missing [ServiceName ServiceVersion]"), err.Error())
}

func TestSLSSpecCreateDirectoryStructureFail(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	spec := slsspec.New()
	err = spec.CreateDirectoryStructure(tmp, specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}, false)
	require.Error(t, err)

	assert.Regexp(t, regexp.MustCompile(`.+ is not a path to foo-1.0.0`), err.Error())
}

func TestSLSSpecCreateDirectoryStructure(t *testing.T) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	rootDir := path.Join(tmp, "foo-1.0.0")
	err = os.Mkdir(rootDir, 0755)
	require.NoError(t, err)

	spec := slsspec.New()
	values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
	err = spec.CreateDirectoryStructure(rootDir, values, false)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(rootDir, "deployment", "manifest.yml"), []byte("test"), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "init.sh"), []byte("test"), 0644)
	require.NoError(t, err)

	specDir, err := specdir.New(rootDir, spec, values, specdir.Validate)
	require.NoError(t, err)

	assert.Equal(t, path.Join(rootDir, "deployment"), specDir.Path(slsspec.Deployment))
	assert.Equal(t, path.Join(rootDir, "deployment", "manifest.yml"), specDir.Path(slsspec.Manifest))
	assert.Equal(t, path.Join(rootDir, "service"), specDir.Path(slsspec.Service))
	assert.Equal(t, path.Join(rootDir, "service", "bin"), specDir.Path(slsspec.ServiceBin))
	assert.Equal(t, path.Join(rootDir, "service", "bin", "init.sh"), specDir.Path(slsspec.InitSh))
}

func TestValidate(t *testing.T) {
	const invalidYML = `[}`

	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name      string
		createDir func(rootDir string)
		skipYML   matcher.Matcher
	}{
		{
			name: "valid structure",
			createDir: func(rootDir string) {
				spec := slsspec.New()
				values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
				err = spec.CreateDirectoryStructure(rootDir, values, false)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "init.sh"), []byte("test"), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "deployment", "manifest.yml"), []byte("key: value"), 0644)
				require.NoError(t, err)
			},
		},
		{
			name: "complex YML is valid",
			createDir: func(rootDir string) {
				spec := slsspec.New()
				values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
				err = spec.CreateDirectoryStructure(rootDir, values, false)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "init.sh"), []byte("test"), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "deployment", "manifest.yml"), []byte(`
# Comments in YAML look like this.

################
# SCALAR TYPES #
################

# Our root object (which continues for the entire document) will be a map,
# which is equivalent to a dictionary, hash or object in other languages.
key: value
another_key: Another value goes here.
a_number_value: 100
scientific_notation: 1e+12
# The number 1 will be interpreted as a number, not a boolean. if you want
# it to be intepreted as a boolean, use true
boolean: true
null_value: null
key with spaces: value
# Notice that strings don't need to be quoted. However, they can be.
however: "A string, enclosed in quotes."
"Keys can be quoted too.": "Useful if you want to put a ':' in your key."

# Multiple-line strings can be written either as a 'literal block' (using |),
# or a 'folded block' (using '>').
literal_block: |
    This entire block of text will be the value of the 'literal_block' key,
    with line breaks being preserved.

    The literal continues until de-dented, and the leading indentation is
    stripped.

        Any lines that are 'more-indented' keep the rest of their indentation -
        these lines will be indented by 4 spaces.
folded_style: >
    This entire block of text will be the value of 'folded_style', but this
    time, all newlines will be replaced with a single space.

    Blank lines, like above, are converted to a newline character.

        'More-indented' lines keep their newlines, too -
        this text will appear over two lines.

####################
# COLLECTION TYPES #
####################

# Nesting is achieved by indentation.
a_nested_map:
    key: value
    another_key: Another Value
    another_nested_map:
        hello: hello

# Maps don't have to have string keys.
0.25: a float key

# Keys can also be complex, like multi-line objects
# We use ? followed by a space to indicate the start of a complex key.
? |
    This is a key
    that has multiple lines
: and this is its value

# Sequences (equivalent to lists or arrays) look like this:
a_sequence:
    - Item 1
    - Item 2
    - 0.5 # sequences can contain disparate types.
    - Item 4
    - key: value
      another_key: another_value
    -
        - This is a sequence
        - inside another sequence

# Since YAML is a superset of JSON, you can also write JSON-style maps and
# sequences:
json_map: {"key": "value"}
json_seq: [3, 2, 1, "takeoff"]

#######################
# EXTRA YAML FEATURES #
#######################

# YAML also has a handy feature called 'anchors', which let you easily duplicate
# content across your document. Both of these keys will have the same value:
anchored_content: &anchor_name This string will appear as the value of two keys.
other_anchor: *anchor_name

# Anchors can be used to duplicate/inherit properties
base: &base
    name: Everyone has same name

foo: &foo
    <<: *base
    age: 10

bar: &bar
    <<: *base
    age: 20

# foo and bar would also have name: Everyone has same name

# YAML also has tags, which you can use to explicitly declare types.
explicit_string: !!str 0.5
# Some parsers implement language specific tags, like this one for Python's
# complex number type.
python_complex_number: !!python/complex 1+2j

####################
# EXTRA YAML TYPES #
####################

# Strings and numbers aren't the only scalars that YAML can understand.
# ISO-formatted date and datetime literals are also parsed.
datetime: 2001-12-15T02:59:43.1Z
datetime_with_spaces: 2001-12-14 21:59:43.10 -5
date: 2002-12-14

# The !!binary tag indicates that a string is actually a base64-encoded
# representation of a binary blob.
gif_file: !!binary |
    R0lGODlhDAAMAIQAAP//9/X17unp5WZmZgAAAOfn515eXvPz7Y6OjuDg4J+fn5
    OTk6enp56enmlpaWNjY6Ojo4SEhP/++f/++f/++f/++f/++f/++f/++f/++f/+
    +f/++f/++f/++f/++f/++SH+Dk1hZGUgd2l0aCBHSU1QACwAAAAADAAMAAAFLC
    AgjoEwnuNAFOhpEMTRiggcz4BNJHrv/zCFcLiwMWYNG84BwwEeECcgggoBADs=

# YAML also has a set type, which looks like this:
set:
    ? item1
    ? item2
    ? item3

# Like Python, sets are just maps with null values; the above is equivalent to:
set2:
    item1: null
    item2: null
    item3: null`), 0644)
				require.NoError(t, err)
			},
		},
		{
			name: "invalid yml file OK if excluded by matcher",
			createDir: func(rootDir string) {
				spec := slsspec.New()
				values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
				err = spec.CreateDirectoryStructure(rootDir, values, false)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "init.sh"), []byte("test"), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "deployment", "manifest.yml"), []byte("key: value"), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "invalid.yaml"), []byte(invalidYML), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "valid.yml"), []byte("key: value"), 0644)
				require.NoError(t, err)
			},
			skipYML: matcher.Path("service"),
		},
	} {
		currTmp, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		rootDir := path.Join(currTmp, "foo-1.0.0")
		err = os.Mkdir(rootDir, 0755)
		require.NoError(t, err)

		currCase.createDir(rootDir)
		err = slsspec.Validate(rootDir, slsspec.TemplateValues("foo", "1.0.0"), currCase.skipYML)
		assert.NoError(t, err, "Case %d: %s", i, currCase.name)
	}
}

func TestValidateFail(t *testing.T) {
	const invalidYML = `[}`

	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		name      string
		createDir func(rootDir string)
		want      string
	}{
		{
			name: "missing manifest.yml",
			createDir: func(rootDir string) {
				spec := slsspec.New()
				values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
				err = spec.CreateDirectoryStructure(rootDir, values, false)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "init.sh"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			want: "foo-1.0.0/deployment/manifest.yml does not exist",
		},
		{
			name: "missing init.sh",
			createDir: func(rootDir string) {
				spec := slsspec.New()
				values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
				err = spec.CreateDirectoryStructure(rootDir, values, false)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "deployment", "manifest.yml"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			want: "foo-1.0.0/service/bin/init.sh does not exist",
		},
		{
			name: "invalid manifest.yml",
			createDir: func(rootDir string) {
				spec := slsspec.New()
				values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
				err = spec.CreateDirectoryStructure(rootDir, values, false)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "init.sh"), []byte("test"), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "deployment", "manifest.yml"), []byte(invalidYML), 0644)
				require.NoError(t, err)
			},
			want: "invalid YML files: [foo-1.0.0/deployment/manifest.yml]\nIf these files are known to be correct, exclude them from validation using the SLS YML validation exclude matcher.",
		},
		{
			name: "multiple invalid yml files",
			createDir: func(rootDir string) {
				spec := slsspec.New()
				values := specdir.TemplateValues{"ServiceName": "foo", "ServiceVersion": "1.0.0"}
				err = spec.CreateDirectoryStructure(rootDir, values, false)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "service", "bin", "init.sh"), []byte("test"), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "deployment", "manifest.yml"), []byte(invalidYML), 0644)
				require.NoError(t, err)
				err = os.MkdirAll(path.Join(rootDir, "var"), 0755)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "var", "invalid.yaml"), []byte(invalidYML), 0644)
				require.NoError(t, err)
				err = ioutil.WriteFile(path.Join(rootDir, "var", "valid.yml"), []byte("key: value"), 0644)
				require.NoError(t, err)
			},
			want: "invalid YML files: [foo-1.0.0/deployment/manifest.yml foo-1.0.0/var/invalid.yaml]\nIf these files are known to be correct, exclude them from validation using the SLS YML validation exclude matcher.",
		},
	} {
		currTmp, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		rootDir := path.Join(currTmp, "foo-1.0.0")
		err = os.Mkdir(rootDir, 0755)
		require.NoError(t, err)

		currCase.createDir(rootDir)
		err = slsspec.Validate(rootDir, slsspec.TemplateValues("foo", "1.0.0"), nil)
		require.Error(t, err, fmt.Sprintf("Case %d: %s", i, currCase.name))

		assert.EqualError(t, err, currCase.want, "Case %d: %s", i, currCase.name)
	}
}
