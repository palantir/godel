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

// Package errorstringer provides functions for converting errors into a a human-readable string representation.
package errorstringer

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type causer interface {
	Cause() error
}
type stackTracer interface {
	StackTrace() errors.StackTrace
}
type errWithStack struct {
	err   error
	msg   string
	stack errors.StackTrace
}

// SingleStack prints the representation of the provided error in the form of all of its messages followed by a single
// stack trace if the error supports such a representation. If the error implements causer and consists of stackTracer
// errors that are all part of the same stack trace, the string will contain all of the messages followed by the longest
// stack trace. If the error cannot be represented in this manner, the returned string is the provided error's "%+v"
// string representation.
func SingleStack(err error) string {
	var stackErrs []errWithStack
	errCause := err
	for errCause != nil {
		if s, ok := errCause.(stackTracer); ok {
			stackErrs = append(stackErrs, errWithStack{
				err:   errCause,
				msg:   errCause.Error(),
				stack: s.StackTrace(),
			})
		}
		if c, ok := errCause.(causer); ok {
			errCause = c.Cause()
		} else {
			break
		}
	}

	if len(stackErrs) == 0 {
		return fmt.Sprintf("%+v", err)
	}

	singleStack := true
	for i := len(stackErrs) - 1; i > 0; i-- {
		if !hasSuffix(stackErrs[i].stack, stackErrs[i-1].stack) {
			singleStack = false
			break
		}
	}

	if !singleStack {
		return fmt.Sprintf("%+v", err)
	}

	for i := 0; i < len(stackErrs)-1; i++ {
		stackErrs[i].msg = currErrMsg(stackErrs[i].err, stackErrs[i+1].err)
	}

	var errs []string
	if errCause != nil {
		// if root cause is non-nil, print its error message if it differs from cause
		stackErrs[len(stackErrs)-1].msg = currErrMsg(stackErrs[len(stackErrs)-1].err, errCause)
		rootErr := errCause.Error()
		if rootErr != stackErrs[len(stackErrs)-1].msg {
			errs = append(errs, errCause.Error())
		}
	}
	for i := len(stackErrs) - 1; i >= 0; i-- {
		errs = append(errs, fmt.Sprintf("%v", stackErrs[i].msg))
	}
	return strings.Join(errs, "\n") + fmt.Sprintf("%+v", stackErrs[len(stackErrs)-1].stack)
}

// StackWithInterleavedMessages prints the representation of the provided error as a stack trace with messages
// interleaved in the relevant locations. If the error implements causer and the error and its causes alternate between
// stackTracer and non-stackTracer errors, the returned string representation is one in which the messages are
// interleaved at the proper positions in the stack and common stack frames in consecutive elements are removed. If the
// provided error is not of this form, the result of printing the provided error using the "%+v" formatting directive is
// returned.
func StackWithInterleavedMessages(err error) string {
	// error is expected to alternate between stackTracer and non-stackTracer errors
	expectStackTracer := true
	var stackErrs []errWithStack
	errCause := err
	for errCause != nil {
		s, isStackTracer := errCause.(stackTracer)
		if isStackTracer {
			stackErrs = append(stackErrs, errWithStack{
				err:   errCause,
				msg:   errCause.Error(),
				stack: s.StackTrace(),
			})
		}
		if c, ok := errCause.(causer); ok {
			if isStackTracer != expectStackTracer {
				// if current error is a causer and does not meet the expectation of whether or not it
				// should be a stackTracer, error is not supported.
				stackErrs = nil
				break
			}
			errCause = c.Cause()
		} else {
			break
		}
		expectStackTracer = !expectStackTracer
	}

	if len(stackErrs) == 0 {
		return fmt.Sprintf("%+v", err)
	}

	for i := len(stackErrs) - 1; i > 0; i-- {
		if hasSuffix(stackErrs[i].stack, stackErrs[i-1].stack) {
			// if the inner stack has the outer stack as a suffix, trim the outer stack from the inner stack
			stackErrs[i].stack = stackErrs[i].stack[0 : len(stackErrs[i].stack)-len(stackErrs[i-1].stack)]
		}
	}

	for i := 0; i < len(stackErrs)-1; i++ {
		stackErrs[i].msg = currErrMsg(stackErrs[i].err, stackErrs[i+1].err)
	}

	var errs []string
	if errCause != nil {
		// if root cause is non-nil, print its error message if it differs from cause
		stackErrs[len(stackErrs)-1].msg = currErrMsg(stackErrs[len(stackErrs)-1].err, errCause)
		rootErr := errCause.Error()
		if rootErr != stackErrs[len(stackErrs)-1].msg {
			errs = append(errs, errCause.Error())
		}
	}
	for i := len(stackErrs) - 1; i >= 0; i-- {
		stack := strings.Replace(fmt.Sprintf("%+v", stackErrs[i].stack), "\n", "\n"+strings.Repeat("\t", 1), -1)
		errs = append(errs, stackErrs[i].msg+stack)
	}
	return strings.Join(errs, "\n")
}

// currErrMsg Returns the string for the current error. If curr.Error() has the suffix ": {{next.Error()}}", returns the
// content of "curr.Error()" up to that suffix.
func currErrMsg(curr error, next error) string {
	msg := curr.Error()
	if idx := strings.LastIndex(msg, ": "+next.Error()); idx != -1 {
		return msg[:idx]
	}
	return msg
}

// hasSuffix returns true if the inner stack trace ends with the outer stack trace, false otherwise.
func hasSuffix(inner errors.StackTrace, outer errors.StackTrace) bool {
	outerIndex := len(outer) - 1
	innerIndex := len(inner) - 1
	for outerIndex >= 0 && innerIndex >= 0 {
		if outer[outerIndex] != inner[innerIndex] {
			break
		}
		outerIndex--
		innerIndex--
	}
	return outerIndex == 0 && innerIndex >= 0
}
