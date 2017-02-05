package pkg2

import ejson "encoding/json"

// Errcheck function.
func Errcheck() {
	helper := func() error {
		return nil
	}
	// does not check for returned error
	helper()
}

// Outparamcheck function.
func Outparamcheck() {
	_ = ejson.Unmarshal(nil, "")
}

// Vet function.
func Vet() string {
	var foo string
	// self-assignment
	foo = foo
	return foo
}

func deadcode() {
	// unused unexported function
}

// Ineffassign function.
func Ineffassign() {
	var kvs [][]string
	kvs = append(kvs, nil)
	// assignment has no result
	kvs = nil
}

// unexported unused constant
const varcheck = "foo"

// Unconvert function.
func Unconvert() string {
	foo := "bar"
	foo = string(foo)
	return foo
}

// This comment will be flagged as bad by Lint.
func Lint() {
}
