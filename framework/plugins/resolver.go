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

package plugins

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/pkg/osarch"
	"github.com/palantir/godel/godelgetter"
)

type resolver interface {
	Resolve(locator locatorWithChecksumsParam, osArch osarch.OSArch, dst string, stdout io.Writer) error
}

func newTemplateResolver(tmpl string) (resolver, error) {
	parsed, err := template.New("resolver").Funcs(funcMap(locatorWithChecksumsParam{}, osarch.OSArch{})).Parse(tmpl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create resolver from template %q", tmpl)
	}
	return &goTemplateResolver{
		tmpl:    parsed,
		tmplSrc: tmpl,
	}, nil
}

type goTemplateResolver struct {
	tmpl    *template.Template
	tmplSrc string
}

func (r goTemplateResolver) Resolve(locator locatorWithChecksumsParam, osArch osarch.OSArch, dst string, stdout io.Writer) error {
	buf := &bytes.Buffer{}
	if err := r.tmpl.Funcs(funcMap(locator, osArch)).Execute(buf, nil); err != nil {
		return errors.Wrapf(err, "failed to execute template %q", r.tmplSrc)
	}
	srcURL := buf.String()

	if _, err := godelgetter.Download(godelgetter.NewPkgSrc(srcURL, ""), dst, stdout); err != nil {
		return errors.Wrapf(err, "failed to resolve artifact at %s", srcURL)
	}
	return nil
}

func funcMap(locator locatorWithChecksumsParam, osArch osarch.OSArch) template.FuncMap {
	return template.FuncMap{
		"Group": func() string {
			return locator.Group
		},
		"GroupPath": func() string {
			return strings.Replace(locator.Group, ".", "/", -1)
		},
		"Product": func() string {
			return locator.Product
		},
		"Version": func() string {
			return locator.Version
		},
		"OS": func() string {
			return osArch.OS
		},
		"Arch": func() string {
			return osArch.Arch
		},
	}
}

// Retrieves the archive for the plugin specified by the provided locator and OSArch and writes it to the provided
// destination path. If the provided locator specifies a resolver, it will be used to retrieve the archive; otherwise,
// the default resolvers will be used in order. If the locator specifies a checksum for the provided OSArch, then it
// will be used to verify the contents of the downloaded archive (note that the checksum is for the file in the archive
// and not the archive itself). Returns an error if the plugin could not be resolved using the resolvers or if a
// checksum was provided and did not match. Note that, if the function resolves an artifact to the destination, the
// artifact will not be removed even if the function returns an error (for example, due to checksums not matching).
func resolvePluginTGZ(locatorWithResolver locatorWithResolverParam, defaultResolvers []resolver, osArch osarch.OSArch, dst string, stdout io.Writer) error {
	const errIndentSpaces = 4

	resolversToUse := defaultResolvers
	if locatorWithResolver.Resolver != nil {
		resolversToUse = []resolver{locatorWithResolver.Resolver}
	}

	success := false
	var errs []string
	for _, resolver := range resolversToUse {
		if err := resolver.Resolve(locatorWithResolver.LocatorWithChecksums, osArch, dst, stdout); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		success = true
		break
	}

	if !success {
		parts := append([]string{fmt.Sprintf("failed to resolve artifact %+v using resolvers:", locatorWithResolver.LocatorWithChecksums)}, errs...)
		return errors.New(strings.Join(parts, fmt.Sprintf("\n%s", strings.Repeat(" ", errIndentSpaces))))
	}

	gotChecksum, err := pluginTGZFileContentHash(dst)
	if err != nil {
		return errors.Wrapf(err, "archive at %s was not structured properly", dst)
	}

	wantChecksum, ok := locatorWithResolver.LocatorWithChecksums.Checksums[osArch]
	if !ok {
		// no checksum present
		return nil
	}

	if wantChecksum != gotChecksum {
		return errors.Errorf("checksum for plugin in archive %s did not match: want %s, got %s", dst, wantChecksum, gotChecksum)
	}
	return nil
}

func pluginTGZFileContentHash(tgzPath string) (string, error) {
	f, err := os.Open(tgzPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open %s", tgzPath)
	}
	defer func() {
		// nothing to do if closing file open for reading fails
		_ = f.Close()
	}()
	return pluginTGZContentHash(f)
}

// Computes the SHA-256 hash for the content of the provided reader, which must be a plugin TGZ.
func pluginTGZContentHash(tgzContentReader io.Reader) (string, error) {
	hasher := sha256.New()
	if err := copyPluginTGZContent(hasher, tgzContentReader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// Verifies that the TGZ content provided by the reader consists of a tar archive that contains a single regular file
// and writes the content of that file to the provided writer. Returns an error if the tar archive does not contain a
// single file (if it contains greater or fewer files or contains non-file entries).
func copyPluginTGZContent(dst io.Writer, tgzContentReader io.Reader) error {
	gzf, err := gzip.NewReader(tgzContentReader)
	if err != nil {
		return errors.Wrapf(err, "failed to create reader")
	}

	tarReader := tar.NewReader(gzf)
	numFiles := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrapf(err, "failed to read tar entry")
		}

		numFiles++
		if header.Typeflag != tar.TypeReg {
			continue
		}

		if numFiles != 1 {
			continue
		}

		if _, err := io.Copy(dst, tarReader); err != nil {
			return errors.Wrapf(err, "failed to read tar file entry")
		}
	}
	if numFiles != 1 {
		return errors.Errorf("archive must contain exactly 1 file, but contained %d", numFiles)
	}
	return nil
}

func sha256ChecksumFile(fPath string) (string, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open %s for reading", f)
	}
	defer func() {
		_ = f.Close()
	}()
	return sha256Checksum(f)
}

func sha256Checksum(r io.Reader) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return "", errors.WithStack(err)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
