// Package werror defines an error type that can store safe and unsafe parameters and can wrap other errors.
package werror

import (
	"context"
	"fmt"

	wparams "github.com/palantir/witchcraft-go-params"
)

var _ Werror = (*werror)(nil)

// Error is identical to calling ErrorWithContext with a context that does not have any wparams parameters.
// DEPRECATED: Please use ErrorWithContextParams instead to ensure that all the wparams parameters that are set on the
// context are included in the error.
func Error(msg string, params ...Param) error {
	return newWerror(msg, nil, params...)
}

// ErrorWithContextParams returns a new error with the provided message and parameters. The returned error also includes any
// wparams parameters that are stored in the context.
//
// The message should not contain any formatted parameters -- instead, use the SafeParam* or UnsafeParam* functions
// to create error parameters.
//
// Example:
//
//	password, ok := config["password"]
//	if !ok {
//		return werror.ErrorWithContextParams(ctx, "configuration is missing password")
//	}
func ErrorWithContextParams(ctx context.Context, msg string, params ...Param) error {
	safe, unsafe := wparams.SafeAndUnsafeParamsFromContext(ctx)
	fullParams := []Param{
		SafeParams(safe),
		UnsafeParams(unsafe),
	}
	fullParams = append(fullParams, params...)
	return newWerror(msg, nil, fullParams...)
}

// Wrap is identical to calling WrapWithContextParams with a context that does not have any wparams parameters.
// DEPRECATED: Please use WrapWithContextParams instead to ensure that all the wparams parameters that are set on the
// context are included in the error.
func Wrap(err error, msg string, params ...Param) error {
	if err == nil {
		return nil
	}
	return newWerror(msg, err, params...)
}

// WrapWithContextParams returns a new error with the provided message and stores the provided error as its cause.
// The returned error also includes any wparams parameters that are stored in the context.
//
// The message should not contain any formatted parameters -- instead use the SafeParam* or UnsafeParam* functions
// to create error parameters.
//
// Example:
//
//	users, err := getUser(userID)
//	if err != nil {
//		return werror.WrapWithContextParams(ctx, err, "failed to get user", werror.SafeParam("userId", userID))
//	}
func WrapWithContextParams(ctx context.Context, err error, msg string, params ...Param) error {
	if err == nil {
		return nil
	}
	safe, unsafe := wparams.SafeAndUnsafeParamsFromContext(ctx)
	fullParams := []Param{
		SafeParams(safe),
		UnsafeParams(unsafe),
	}
	fullParams = append(fullParams, params...)
	return newWerror(msg, err, fullParams...)
}

// Convert err to werror error.
//
// If err is not a werror-based error, then a new werror error is created using the message from err.
// Otherwise, returns unchanged err.
//
// Example:
//
//	file, err := os.Open("file.txt")
//	if err != nil {
//		return werror.Convert(err)
//	}
func Convert(err error) error {
	if err == nil {
		return err
	}
	switch err.(type) {
	case Werror:
		return err
	default:
		return newWerror("", err)
	}
}

// RootCause returns the initial cause of an error.
//
// Traverses the cause hierarchy until it reaches an error which has no cause and returns that error.
func RootCause(err error) error {
	for {
		causer, ok := err.(Causer)
		if !ok {
			return err
		}
		cause := causer.Cause()
		if cause == nil {
			return err
		}
		err = cause
	}
}

// ParamsFromError returns all of the safe and unsafe parameters stored in the provided error.
//
// If the error implements the Causer interface, then the returned parameters will include all of the parameters stored
// in the causes as well.
//
// All of the keys and parameters of the map are flattened.
//
// Parameters are added from the outermost error to the innermost error. This means that, if multiple errors declare
// different values for the same keys, the values for the most specific (deepest) error will be the ones in the returned
// maps.
func ParamsFromError(err error) (safeParams map[string]interface{}, unsafeParams map[string]interface{}) {
	safeParams = make(map[string]interface{})
	unsafeParams = make(map[string]interface{})
	if err != nil {
		visitErrorParams(err, func(k string, v interface{}, safe bool) {
			if safe {
				safeParams[k] = v
			} else {
				unsafeParams[k] = v
			}
		})
	}
	return safeParams, unsafeParams
}

// ParamFromError returns the value of the parameter for the given key, or nil if no such key exists. Checks the
// parameters of the provided error and all of its causes. If the error and its causes contain multiple values for the
// same key, the most specific (deepest) value will be returned.
func ParamFromError(err error, key string) (value interface{}, safe bool) {
	visitErrorParams(err, func(k string, v interface{}, s bool) {
		if k == key {
			value = v
			safe = s
		}
	})
	return value, safe
}

// visitErrorParams calls the provided visitor function on all of the parameters stored in the provided error and any of
// its causes. The function is invoked on all of the parameters stored in the provided error, then all of the parameters
// in the cause of the provided error, and so on. There are no guarantees made about the order in which the parameters
// will be called for a given error.
func visitErrorParams(err error, visitor func(k string, v interface{}, safe bool)) {
	allErrs := []error{err}
	for currErr := err; ; {
		causer, ok := currErr.(Causer)
		if !ok || causer.Cause() == nil {
			// current error does not have a cause
			break
		}
		allErrs = append(allErrs, causer.Cause())
		currErr = causer.Cause()
	}
	for _, currErr := range allErrs {
		if ps, ok := currErr.(wparams.ParamStorer); ok {
			for k, v := range ps.SafeParams() {
				visitor(k, v, true)
			}
			for k, v := range ps.UnsafeParams() {
				visitor(k, v, false)
			}
		}
	}
}

// Werror is an error type consisting of an underlying error, stacktrace, underlying causes, and safe and unsafe
// params associated with that error.
type Werror interface {
	error
	fmt.Formatter
	Causer
	StackTracer
	wparams.ParamStorer

	Message() string
}

// werror is an error type consisting of an underlying error and safe and unsafe params associated with that error.
type werror struct {
	message string
	cause   error
	stack   StackTrace
	params  map[string]paramValue
}

type paramValue struct {
	safe  bool
	value interface{}
}

// Causer interface is compatible with the interface used by pkg/errors.
type Causer interface {
	Cause() error
}

func newWerror(message string, cause error, params ...Param) error {
	we := &werror{
		message: message,
		cause:   cause,
		stack:   NewStackTraceWithSkip(1),
		params:  make(map[string]paramValue),
	}
	for _, p := range params {
		p.apply(we)
	}
	return we
}

// Error returns the message for this error by delegating to the stored error. The error consists only of the message
// and does not include any other information such as safe/unsafe parameters or cause.
func (e *werror) Error() string {
	if e.cause == nil {
		return e.message
	}
	if e.message == "" {
		return e.cause.Error()
	}
	return e.message + ": " + e.cause.Error()
}

// Cause returns the underlying cause of this error or nil if there is none.
func (e *werror) Cause() error {
	return e.cause
}

// Unwrap returns the wrapped error. Exists to support Go error Is/As functions introduced in Go 1.13.
func (e *werror) Unwrap() error {
	return e.cause
}

// StackTrace returns the Stacktracer for this error or nil if there is none.
func (e *werror) StackTrace() StackTrace {
	return e.stack
}

// Message returns the message string for this error.
func (e *werror) Message() string {
	return e.message
}

// SafeParams returns params from this error and any underlying causes. If the error and its causes
// contain multiple values for the same key, the most specific (deepest) value will be returned.
func (e *werror) SafeParams() map[string]interface{} {
	safe, _ := ParamsFromError(e.cause)
	for k, v := range e.params {
		if v.safe {
			if _, exists := safe[k]; !exists {
				safe[k] = v.value
			}
		}
	}
	return safe
}

// UnsafeParams returns params from this error and any underlying causes. If the error and its causes
// contain multiple values for the same key, the most specific (deepest) value will be returned.
func (e *werror) UnsafeParams() map[string]interface{} {
	_, unsafe := ParamsFromError(e.cause)
	for k, v := range e.params {
		if !v.safe {
			if _, exists := unsafe[k]; !exists {
				unsafe[k] = v.value
			}
		}
	}
	return unsafe
}

// Format formats the error using the provided format state. Delegates to stored error.
func (e *werror) Format(state fmt.State, verb rune) {
	safe := make(map[string]interface{})
	for k, v := range e.params {
		if v.safe {
			safe[k] = v.value
		}
	}
	Format(e, safe, state, verb)
}

// Format formats a Werror using the provided format state. This is a utility method that can
// be used by other implementations of Werror. The safeParams argument is expected to include
// safe params for this error only, not for any underlying causes.
func Format(err Werror, safeParams map[string]interface{}, state fmt.State, verb rune) {
	if verb == 'v' && state.Flag('+') {
		// Multi-line extra verbose format starts with cause first followed up by current error metadata.
		formatCause(err, state, verb)
		formatMessage(err, state, verb)
		formatParameters(err, safeParams, state, verb)
		formatStack(err, state, verb)
	} else {
		formatMessage(err, state, verb)
		formatParameters(err, safeParams, state, verb)
		formatStack(err, state, verb)
		formatCause(err, state, verb)
	}
}

func formatMessage(err Werror, state fmt.State, verb rune) {
	if err.Message() == "" {
		return
	}
	switch verb {
	case 's', 'q', 'v':
		_, _ = fmt.Fprint(state, err.Message())
	}
}

func formatParameters(err Werror, safeParams map[string]interface{}, state fmt.State, verb rune) {
	if len(safeParams) == 0 {
		return
	}
	if verb != 'v' {
		return
	}
	if err.Message() != "" {
		// Whitespace before the message.
		_, _ = fmt.Fprint(state, " ")
	}
	_, _ = fmt.Fprintf(state, "%+v", safeParams)
}

func formatStack(err Werror, state fmt.State, verb rune) {
	if err.StackTrace() == nil {
		return
	}
	if verb != 'v' || !state.Flag('+') {
		return
	}
	err.StackTrace().Format(state, verb)
}

func formatCause(err Werror, state fmt.State, verb rune) {
	if err.Cause() == nil {
		return
	}
	var prefix string
	if err.Message() != "" || (verb == 'v' && len(err.SafeParams()) > 0) {
		prefix = ": "
	}
	switch verb {
	case 'v':
		if state.Flag('+') {
			_, _ = fmt.Fprintf(state, "%+v\n", err.Cause())
		} else {
			_, _ = fmt.Fprintf(state, "%s%v", prefix, err.Cause())
		}
	case 's':
		_, _ = fmt.Fprintf(state, "%s%s", prefix, err.Cause())
	case 'q':
		_, _ = fmt.Fprintf(state, "%s%q", prefix, err.Cause())
	}
}
