// MIT License
//
// Copyright (c) 2016 Nick Miyake
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package errorstringer_test

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/nmiyake/pkg/errorstringer"
)

func TestSingleStack(t *testing.T) {
	cases := []struct {
		err  error
		want []string
	}{
		// regular error prints directly
		{io.EOF, []string{"EOF"}},
		// multiple errors wrapped with "Wrap" coalesce into a single stack
		{loadConfig(), []string{
			"EOF",
			"failed to open foo",
			"failed to parse file",
			"failed to load config",
			"github.com/nmiyake/pkg/errorstringer_test.openFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.parseFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfig",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestSingleStack",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
		}},
		// if error contains multiple stacks that don't coalesce, prints "%+v" output
		{loadConfigInChannel(), []string{
			"EOF",
			"failed to open foo in channel",
			"github.com/nmiyake/pkg/errorstringer_test.openFileInChannel.func1",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"github.com/nmiyake/pkg/errorstringer_test.parseFileInChannel",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfigInChannel",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestSingleStack",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/src/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/src/runtime/asm_amd64.s:[0-9]+",
			"failed to load file",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfigInChannel",
			"\t.+/src/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestSingleStack",
			"\t.+/src/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/src/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/src/runtime/asm_amd64.s:[0-9]+",
		}},
	}

	for i, currCase := range cases {
		got := strings.Split(errorstringer.SingleStack(currCase.err), "\n")

		if len(got) != len(currCase.want) {
			t.Errorf("Case %d:\nwant: %v\ngot:  %v", i, currCase.want, got)
			continue
		}

		for j := range got {
			if !regexp.MustCompile("^" + currCase.want[j] + "$").MatchString(got[j]) {
				t.Errorf("Case %d, line %d:\nwant: %v\ngot:  %v", i, j, currCase.want[j], got[j])
				break
			}
		}
	}
}

func TestStackWithInterleavedMessages(t *testing.T) {
	cases := []struct {
		err  error
		want []string
	}{
		// regular error prints directly
		{io.EOF, []string{"EOF"}},
		// multiple errors wrapped with "Wrap" annotates stack with messages in correct spots
		{loadConfig(), []string{
			"EOF",
			"failed to open foo",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.openFile",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"failed to parse file",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.parseFile",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"failed to load config",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.loadConfig",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.TestStackWithInterleavedMessages",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"\ttesting.tRunner",
			"\t\t.+/testing/testing.go:[0-9]+",
			"\truntime.goexit",
			"\t\t.+/runtime/asm_amd64.s:[0-9]+",
		}},
		// if error causes do not alternate between stackTracer and non-stackTracer, prints "%+v" output
		{loadConfigWithMsg(), []string{
			"EOF",
			"failed to open foo",
			"github.com/nmiyake/pkg/errorstringer_test.openFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.parseFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfigWithMsg",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestStackWithInterleavedMessages",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"github.com/nmiyake/pkg/errorstringer_test.parseFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfigWithMsg",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestStackWithInterleavedMessages",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to load config",
		}},
		// if error causes do not alternate between stackTracer and non-stackTracer, prints "%+v" output
		{loadConfigWithStack(), []string{
			"EOF",
			"failed to open foo",
			"github.com/nmiyake/pkg/errorstringer_test.openFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.parseFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfigWithStack",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestStackWithInterleavedMessages",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"github.com/nmiyake/pkg/errorstringer_test.parseFile",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfigWithStack",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestStackWithInterleavedMessages",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.loadConfigWithStack",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"github.com/nmiyake/pkg/errorstringer_test.TestStackWithInterleavedMessages",
			"\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"testing.tRunner",
			"\t.+/testing/testing.go:[0-9]+",
			"runtime.goexit",
			"\t.+/runtime/asm_amd64.s:[0-9]+",
		}},
		// if error contains multiple stacks, all stacks are included and annotated at correct location
		{loadConfigInChannel(), []string{
			"EOF",
			"failed to open foo in channel",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.openFileInChannel.func1",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"\truntime.goexit",
			"\t\t.+/runtime/asm_amd64.s:[0-9]+",
			"failed to parse file",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.parseFileInChannel",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"failed to load file",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.loadConfigInChannel",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"\tgithub.com/nmiyake/pkg/errorstringer_test.TestStackWithInterleavedMessages",
			"\t\t.+/github.com/nmiyake/pkg/errorstringer/errorstringer_test.go:[0-9]+",
			"\ttesting.tRunner",
			"\t\t.+/src/testing/testing.go:[0-9]+",
			"\truntime.goexit",
			"\t\t.+/src/runtime/asm_amd64.s:[0-9]+",
		}},
	}

	for i, currCase := range cases {
		got := strings.Split(errorstringer.StackWithInterleavedMessages(currCase.err), "\n")

		if len(got) != len(currCase.want) {
			t.Errorf("Case %d:\nwant: %v\ngot:  %v", i, currCase.want, got)
			continue
		}

		for j := range got {
			if !regexp.MustCompile("^" + currCase.want[j] + "$").MatchString(got[j]) {
				t.Errorf("Case %d, line %d:\nwant: %q\ngot:  %q", i, j, currCase.want[j], got[j])
				break
			}
		}
	}
}

func loadConfig() error {
	return errors.Wrap(parseFile(), "failed to load config")
}

func loadConfigWithMsg() error {
	return errors.WithMessage(parseFile(), "failed to load config")
}

func loadConfigWithStack() error {
	return errors.WithStack(parseFile())
}

func parseFile() error {
	return errors.Wrap(openFile(), "failed to parse file")
}

func openFile() error {
	return errors.Wrap(io.EOF, "failed to open foo")
}

func loadConfigInChannel() error {
	return errors.Wrapf(parseFileInChannel(), "failed to load file")
}

func parseFileInChannel() error {
	return errors.Wrap(openFileInChannel(), "failed to parse file")
}

func openFileInChannel() error {
	errs := make(chan error)
	go func() { errs <- errors.Wrap(io.EOF, "failed to open foo in channel") }()
	return <-errs
}
