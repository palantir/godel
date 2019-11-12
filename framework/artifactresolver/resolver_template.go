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

package artifactresolver

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/palantir/godel/v2/godelgetter"
	"github.com/palantir/godel/v2/pkg/osarch"
)

func NewTemplateResolver(tmpl string) (Resolver, error) {
	parsed, err := template.New("resolver").Funcs(funcMap(LocatorParam{}, osarch.OSArch{})).Parse(tmpl)
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

func (r goTemplateResolver) Resolve(locator LocatorParam, osArch osarch.OSArch, dst string, stdout io.Writer) error {
	buf := &bytes.Buffer{}
	if err := r.tmpl.Funcs(funcMap(locator, osArch)).Execute(buf, nil); err != nil {
		return errors.Wrapf(err, "failed to execute template %q", r.tmplSrc)
	}
	srcURL := buf.String()

	if err := godelgetter.Download(godelgetter.NewPkgSrc(srcURL, ""), dst, stdout); err != nil {
		return errors.Wrapf(err, "failed to resolve artifact at %s", srcURL)
	}
	return nil
}

func funcMap(locator LocatorParam, osArch osarch.OSArch) template.FuncMap {
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
