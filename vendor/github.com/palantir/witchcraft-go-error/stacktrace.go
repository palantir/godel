package werror

import (
	"fmt"
	"runtime"

	"github.com/palantir/witchcraft-go-error/internal/errors"
)

var _ StackTrace = (*stack)(nil)

// StackTrace provides formatting for an underlying stack trace.
type StackTrace interface {
	fmt.Formatter
}

// StackTracer provides the behavior necessary to retrieve a StackTrace formatter.
type StackTracer interface {
	StackTrace() StackTrace
}

// NewStackTrace creates a new StackTrace, constructed by collecting program counters from runtime callers.
func NewStackTrace() StackTrace {
	return NewStackTraceWithSkip(1)
}

// NewStackTraceWithSkip creates a new StackTrace that skips an additional `skip` stack frames.
func NewStackTraceWithSkip(skip int) StackTrace {
	const depth = 32
	var pcs [depth]uintptr
	// Changing this back to "3" by default. Most callers have only a single level of indirection. For newWerror
	// specifically, which is always called indirectly, we now call this with skip of "1".
	n := runtime.Callers(skip+3, pcs[:])
	var st stack = pcs[0:n]
	return &st
}

// stack represents a stack of program counters.
type stack []uintptr

func (s *stack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case state.Flag('+'):
			for _, pc := range *s {
				f := errors.Frame(pc)
				_, _ = fmt.Fprintf(state, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*s)[i])
	}
	return f
}
