<p align="right">
<a href="https://autorelease.general.dmz.palantir.tech/palantir/witchcraft-go-params"><img src="https://img.shields.io/badge/Perform%20an-Autorelease-success.svg" alt="Autorelease"></a>
</p>

witchcraft-go-params
====================
[![](https://godoc.org/github.com/palantir/witchcraft-go-params?status.svg)](http://godoc.org/github.com/palantir/witchcraft-go-params)

`witchcraft-go-params` defines the `wparams` package, which provides the `ParamStorer` interface and functions for 
storing and retrieving `ParamStorer` implementations in from a context.

Conceptually, "params" are values that are associated with a specific key that provide context about an operation that
is being performed. Params are categorized as "safe" or "unsafe" -- "safe" params are parameters which are considered
safe to ship/export/expose off-premises, while "unsafe" parameters are parameters that should not leave the premises.
Param values are typically used by things such as loggers and errors to provide further context for an operation. Keys
are case-sensitive and must be unique across both safe and unsafe parameters.

The following is a short example of a canonical use case:

```go
type UserID int64

func UpdateUserInfo(ctx context.Context, userID UserID, info UserInfo) error {
	ctx = wparams.ContextWithSafeParam(ctx, "userId", userID)
	
	svc1log.FromContext(ctx).Info("Updating user information", svc1log.Params(wparams.ParamsFromContext(ctx)))
	if err := validateInput(ctx, info); err != nil {
		return err
	}
	// ...
	return nil
}

func validateInput(ctx context.Context, info UserInfo) error {
    if info.Name == "" {
        return werror.Error("invalid user info", werror.Params(wparams.ParamsFromContext(ctx)))	
    }
	// ...
	return nil
}
```

In this example, "userId" and its value are registered as a safe parameter on the context at the beginning of the
`UpdateUserInfo` function. This function (and the other functions that it calls) extract the safe and unsafe parameters
from the context when performing operations such as logging or creating errors. This makes it such that the logger
messages and errors are provided with all of the relevant parameters for their operation, along with the knowledge of
whether the parameters are safe and unsafe.

If a specific type will be recorded as a parameter often and there is a sensible default value for its name and 
safe/unsafe status, it may make sense to have the type implement the `ParamStorer` interface directly. For example, if
`UserID` is known to be safe and should always be recorded as "userId", the example above could be updated as follows:

```go
type UserID int64

func (id UserID) SafeParams() map[string]interface{} {
	return map[string]interface{}{
		"userId": id,
	}
}

func (id UserID) UnsafeParams() map[string]interface{} {
	return nil
}

func UpdateUserInfo(ctx context.Context, userID UserID, info UserInfo) error {
	ctx = wparams.ContextWithParams(ctx, userID)
    // the rest is the same as the first example
}
```

This pattern allows a type to dictate its default behavior for its name and whether it is safe/unsafe, which makes it
easier to ensure that the parameter is consistent across various usages.

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).
