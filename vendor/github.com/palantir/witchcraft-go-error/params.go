package werror

import (
	wparams "github.com/palantir/witchcraft-go-params"
)

type Param interface {
	apply(*werror)
}

type param func(*werror)

func (p param) apply(e *werror) {
	p(e)
}

func SafeParam(key string, val any) Param {
	return SafeParams(map[string]any{key: val})
}

func SafeParams(vals map[string]any) Param {
	return paramsHelper(vals, true)
}

func UnsafeParam(key string, val any) Param {
	return UnsafeParams(map[string]any{key: val})
}

func UnsafeParams(vals map[string]any) Param {
	return paramsHelper(vals, false)
}

func paramsHelper(vals map[string]any, safe bool) Param {
	return param(func(z *werror) {
		for k, v := range vals {
			z.params[k] = paramValue{
				safe:  safe,
				value: v,
			}
		}
	})
}

func SafeAndUnsafeParams(safe, unsafe map[string]any) Param {
	return param(func(z *werror) {
		SafeParams(safe).apply(z)
		UnsafeParams(unsafe).apply(z)
	})
}

func Params(object wparams.ParamStorer) Param {
	return param(func(z *werror) {
		if object != nil {
			SafeParams(object.SafeParams()).apply(z)
			UnsafeParams(object.UnsafeParams()).apply(z)
		}
	})
}
