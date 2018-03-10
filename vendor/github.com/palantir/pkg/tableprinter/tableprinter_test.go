// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tableprinter

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"text/tabwriter"
)

var (
	b        = new(bytes.Buffer)
	tw       = tabwriter.NewWriter(b, 0, 0, 1, ' ', uint(0))
	cgetters = map[string]ColumnGetter{
		"service": wrap(service.Name),
		"host":    wrap(service.Host),
		"product": wrap(service.ProductName),
		"version": wrap(service.ProductVersion),
	}
)

type service struct {
	name string
	host string
	product
}

func (s service) Name() string {
	return s.name
}

func (s service) Host() string {
	return s.host
}

func (s service) ProductName() string {
	return s.product.name
}

func (s service) ProductVersion() string {
	return s.product.version
}

type product struct {
	name    string
	version string
}

func wrap(getter func(service) string) ColumnGetter {
	return func(row interface{}) string {
		return getter(row.(service))
	}
}

var tests = []struct {
	testname      string
	columns       []string
	rows          []interface{}
	includeHeader bool
	expected      string
}{
	{
		"1 service",
		[]string{"service", "host", "product", "version"},
		[]interface{}{
			service{
				"foo",
				"localhost",
				product{
					"nvim",
					"0.1.5",
				},
			},
		},
		true,
		"service host      product version\nfoo     localhost nvim    0.1.5\n",
	},
	{
		"2 services sorted by name",
		[]string{"service", "host", "product", "version"},
		[]interface{}{
			service{
				"bar",
				"localhost",
				product{
					"vim",
					"7.4",
				},
			},
			service{
				"foo",
				"localhost",
				product{
					"nvim",
					"0.1.5",
				},
			},
		},
		true,
		"service host      product version\nbar     localhost vim     7.4\nfoo     localhost nvim    0.1.5\n",
	},
	{
		"2 services sorted by product",
		[]string{"product", "version", "service", "host"},
		[]interface{}{
			service{
				"bar",
				"localhost",
				product{
					"vim",
					"7.4",
				},
			},
			service{
				"foo",
				"localhost",
				product{
					"nvim",
					"0.1.5",
				},
			},
		},
		true,
		"product version service host\nnvim    0.1.5   foo     localhost\nvim     7.4     bar     localhost\n",
	},
	{
		"2 services sorted by product with no header",
		[]string{"product", "version", "service", "host"},
		[]interface{}{
			service{
				"bar",
				"localhost",
				product{
					"vim",
					"7.4",
				},
			},
			service{
				"foo",
				"localhost",
				product{
					"nvim",
					"0.1.5",
				},
			},
		},
		false,
		"nvim 0.1.5 foo localhost\nvim  7.4   bar localhost\n",
	},
}

func Test(t *testing.T) {
	for _, e := range tests {
		check(t, e.testname, e.columns, e.rows, e.includeHeader, e.expected)
	}
}

func TestNonexistentColumn(t *testing.T) {
	b.Truncate(0)
	p := New(tw, cgetters, true, true, false)
	dne := "DNE" // this column isn't defined
	err := p.Print([]string{dne}, nil)
	switch {
	case err == nil:
		t.Fatalf("Print should fail when given a nonexistent column: %q", dne)
	case !strings.Contains(err.Error(), "undefined column(s)"):
		t.Fatalf("Print should inform the caller that there are undefined ColumnGetter(s): %v", err.Error())
	}
}

func TestErrorDuringWrite(t *testing.T) {
	b.Truncate(0)
	p := New(tabwriter.NewWriter(errorWriter{}, 0, 0, 0, ' ', uint(0)), cgetters, true, true, false)
	err := p.Print([]string{"service"}, nil)
	switch {
	case err == nil:
		t.Errorf("Print should fail when write errors")
	case !strings.Contains(err.Error(), "failed to write"):
		t.Fatalf("Print should inform the caller that it failed to write: %v", err.Error())
	}
}

func check(t *testing.T, testname string, columns []string, rows []interface{}, includeHeader bool, expected string) {
	b.Truncate(0)
	p := New(tw, cgetters, true, includeHeader, false)
	if err := p.Print(columns, rows); err != nil {
		t.Fatalf(strings.Join(
			[]string{"--- test: %s", "--- error: %v"}, "\n",
		), testname, err)
	}
	if b.String() != expected {
		t.Errorf(strings.Join(
			[]string{"--- test: %s", "--- columns: %v", "--- rows:%v", "--- found: %q", "--- expected: %q"}, "\n",
		), testname, columns, rows, b.String(), expected)
	}
}

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("cannot write")
}
