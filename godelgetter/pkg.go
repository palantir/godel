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

package godelgetter

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func NewPkgSrc(srcPath, checksum string) PkgSrc {
	pkg := basePkg{
		path:     srcPath,
		checksum: checksum,
	}
	if strings.HasPrefix(srcPath, "http://") || strings.HasPrefix(srcPath, "https://") {
		return &remotePkg{basePkg: pkg}
	}
	return &localFilePkg{basePkg: pkg}
}

type PkgSrc interface {
	// Returns the name of this package.
	Name() string
	// Returns the path to the source of this package.
	Path() string
	// Returns the expected SHA-256 checksum for the package. If this function returns an empty string, then a checksum
	// will not be performed.
	Checksum() string
	// Returns a reader that can be used to read the package and the size of the package. Reader will be open and ready
	// for reads -- the caller is responsible for closing the reader when done.
	Reader() (io.ReadCloser, int64, error)
}

type basePkg struct {
	path     string
	checksum string
}

func (p *basePkg) Name() string {
	return path.Base(p.path)
}

func (p *basePkg) Path() string {
	return p.path
}

func (p *basePkg) Checksum() string {
	return p.checksum
}

type remotePkg struct {
	basePkg
}

func (p *remotePkg) Reader() (io.ReadCloser, int64, error) {
	response, err := http.Get(p.path)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "get call for URL %s failed", p.path)
	}
	if response.StatusCode >= 400 {
		return nil, 0, errors.Errorf("request for URL %s returned status code %d", p.path, response.StatusCode)
	}
	return response.Body, response.ContentLength, nil
}

type localFilePkg struct {
	basePkg
}

func (p *localFilePkg) Reader() (io.ReadCloser, int64, error) {
	localTgzFileInfo, err := os.Stat(p.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, errors.Errorf("%s does not exist", p.path)
		}
		return nil, 0, errors.WithStack(err)
	} else if localTgzFileInfo.IsDir() {
		return nil, 0, errors.Errorf("%s is a directory", p.path)
	}
	srcTgzFile, err := os.Open(p.path)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to open %s", p.path)
	}
	return srcTgzFile, localTgzFileInfo.Size(), nil
}
