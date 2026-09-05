package wparams

import (
	"context"
)

type witchcraftParamsContextKeyType string

const wParamsContextKey = witchcraftParamsContextKeyType("witchcraftParams")

// ContextWithParamStorers returns a copy of the provided context that contains all of the safe and unsafe parameters
// provided by the provided ParamStorers. If the provided context already has safe/unsafe params, the newly returned
// context will contain the result of merging the previous parameters with the provided parameters.
func ContextWithParamStorers(ctx context.Context, params ...ParamStorer) context.Context {
	return context.WithValue(ctx, wParamsContextKey, NewParamStorer(append([]ParamStorer{ParamStorerFromContext(ctx)}, params...)...))
}

// ContextWithSafeParam returns a copy of the provided context that contains the provided safe parameter. If the
// provided context already has safe/unsafe params, the newly returned context will contain the result of merging the
// previous parameters with the provided parameter.
func ContextWithSafeParam(ctx context.Context, key string, value interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewSafeParam(key, value))
}

// ContextWithSafeParams returns a copy of the provided context that contains the provided safe parameters. If the
// provided context already has safe/unsafe params, the newly returned context will contain the result of merging the
// previous parameters with the provided parameters.
func ContextWithSafeParams(ctx context.Context, safeParams map[string]interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewSafeParamStorer(safeParams))
}

// ContextWithUnsafeParam returns a copy of the provided context that contains the provided unsafe parameter. If the
// provided context already has safe/unsafe params, the newly returned context will contain the result of merging the
// previous parameters with the provided parameter.
func ContextWithUnsafeParam(ctx context.Context, key string, value interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewUnsafeParam(key, value))
}

// ContextWithUnsafeParams returns a copy of the provided context that contains the provided unsafe parameters. If the
// provided context already has safe/unsafe params, the newly returned context will contain the result of merging the
// previous parameters with the provided parameters.
func ContextWithUnsafeParams(ctx context.Context, unsafeParams map[string]interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewUnsafeParamStorer(unsafeParams))
}

// ContextWithSafeAndUnsafeParams returns a copy of the provided context that contains the provided safe and unsafe
// parameters. If the provided context already has safe/unsafe params, the newly returned context will contain the
// result of merging the previous parameters with the provided parameters.
func ContextWithSafeAndUnsafeParams(ctx context.Context, safeParams, unsafeParams map[string]interface{}) context.Context {
	return ContextWithParamStorers(ctx, NewSafeAndUnsafeParamStorer(safeParams, unsafeParams))
}

// ParamStorerFromContext returns the ParamStorer stored in the provided context. Returns nil if the provided context
// does not contain a ParamStorer.
func ParamStorerFromContext(ctx context.Context) ParamStorer {
	val := ctx.Value(wParamsContextKey)
	if paramStorer, ok := val.(ParamStorer); ok {
		return paramStorer
	}
	return nil
}

// SafeAndUnsafeParamsFromContext returns the safe and unsafe parameters stored in the ParamStorer returned by
// ParamStorerFromContext for the provided context. Returns nil maps if the provided context does not have a
// ParamStorer.
func SafeAndUnsafeParamsFromContext(ctx context.Context) (safeParams map[string]interface{}, unsafeParams map[string]interface{}) {
	storer := ParamStorerFromContext(ctx)
	if storer == nil {
		return nil, nil
	}
	return storer.SafeParams(), storer.UnsafeParams()
}
