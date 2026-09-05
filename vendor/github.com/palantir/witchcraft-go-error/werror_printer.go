package werror

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

// GenerateErrorString will attempt to pretty print an error depending on its underlying type
// If it is a werror then:
// 1) Each message and params will be groups together on a separate line
// 2) Only the deepest werror stacktrace will be printed
// 3) GenerateErrorString will be called recursively to pretty print underlying errors as well
// If the error implements the fmt.Formatter interface, then it will be printed verbosely
// Otherwise, the error's underlying Error() function will be called and returned
func GenerateErrorString(err error, outputEveryCallingStack bool) string {
	if werror, ok := err.(Werror); ok {
		return generateWerrorString(werror, outputEveryCallingStack)
	}
	if fancy, ok := err.(fmt.Formatter); ok {
		// This is a rich error type, like those produced by github.com/pkg/errors.
		return fmt.Sprintf("%+v", fancy)
	}
	return err.Error()
}

func generateWerrorString(err Werror, outputEveryCallingStack bool) string {
	var buffer bytes.Buffer
	writeMessage(err, &buffer)
	writeParams(err, &buffer)
	writeCause(err, &buffer, outputEveryCallingStack)
	writeStack(err, &buffer, outputEveryCallingStack)
	return buffer.String()
}

func writeMessage(err Werror, buffer *bytes.Buffer) {
	if err.Message() == "" {
		return
	}
	buffer.WriteString(err.Message())
}

func writeParams(err Werror, buffer *bytes.Buffer) {
	safeParams := getSafeParamsAtCurrentLevel(err)
	var safeKeys []string
	for k := range safeParams {
		safeKeys = append(safeKeys, k)
	}
	sort.Strings(safeKeys)
	messageAndParams := err.Message() != "" && len(safeParams) != 0
	messageOrParams := err.Message() != "" || len(safeParams) != 0
	if messageAndParams {
		buffer.WriteString(" ")
	}
	for _, safeKey := range safeKeys {
		safeValue := safeParams[safeKey]
		if v := reflect.ValueOf(safeValue); v.Kind() == reflect.Ptr && !v.IsNil() {
			safeValue = v.Elem().Interface()
		}
		buffer.WriteString(fmt.Sprintf("%+v:%+v", safeKey, safeValue))
		// If it is not the last param, add a separator
		if !(safeKeys[len(safeKeys)-1] == safeKey) {
			buffer.WriteString(", ")
		}
	}
	if messageOrParams {
		buffer.WriteString("\n")
	}
}

func getSafeParamsAtCurrentLevel(err Werror) map[string]interface{} {
	safeParamsAtThisLevel := make(map[string]interface{}, 0)
	childSafeParams := getChildSafeParams(err)
	for k, v := range err.SafeParams() {
		_, ok := childSafeParams[k]
		if ok {
			continue
		}
		safeParamsAtThisLevel[k] = v
	}
	return safeParamsAtThisLevel
}

func getChildSafeParams(err Werror) map[string]interface{} {
	if err.Cause() == nil {
		return make(map[string]interface{}, 0)
	}
	causeAsWerror, ok := err.Cause().(Werror)
	if !ok {
		return make(map[string]interface{}, 0)
	}
	return causeAsWerror.SafeParams()
}

func writeCause(err Werror, buffer *bytes.Buffer, outputEveryCallingStack bool) {
	if err.Cause() != nil {
		buffer.WriteString(GenerateErrorString(err.Cause(), outputEveryCallingStack))
	}
}

func writeStack(err Werror, buffer *bytes.Buffer, outputEveryCallingStack bool) {
	if _, ok := err.Cause().(Werror); ok {
		if !outputEveryCallingStack {
			return
		}
	}
	buffer.WriteString(fmt.Sprintf("%+v", err.StackTrace()))
}
