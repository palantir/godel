package wparams

// ParamStorer is a type that stores safe and unsafe parameters. Keys should be unique across both SafeParams and
// UnsafeParams (that is, if a key occurs in one map, it should not occur in the other). For performance reasons,
// the maps returned by SafeParams and UnsafeParams are references to the underlying storage and should not be modified
// by the caller.
type ParamStorer interface {
	SafeParams() map[string]interface{}
	UnsafeParams() map[string]interface{}
}

// NewParamStorer returns a new ParamStorer that stores all of the params in the provided ParamStorer inputs. The params
// are added from the param storers in the order in which they are provided, and for each individual param storer all of
// the safe params are added before the unsafe params while maintaining key uniqueness across both safe and unsafe
// parameters. This means that, if the same parameter is provided by multiple ParamStorer inputs, the returned
// ParamStorer will have the key (including safe/unsafe type) and value as provided by the last ParamStorer (for
// example, if an unsafe key/value pair is provided by one ParamStorer and a later ParamStorer specifies a safe
// key/value pair with the same key, the returned ParamStorer will store the last safe key/value pair).
func NewParamStorer(paramStorers ...ParamStorer) ParamStorer {
	collector := &mapParamStorer{}
	for _, storer := range paramStorers {
		collector.copyFrom(storer)
	}
	return collector
}

// NewSafeParamStorer returns a new ParamStorer that stores the provided parameters as SafeParams.
func NewSafeParamStorer(safeParams map[string]interface{}) ParamStorer {
	return NewSafeAndUnsafeParamStorer(safeParams, nil)
}

// NewSafeParam returns a new ParamStorer that stores a single safe parameter.
func NewSafeParam(key string, value interface{}) ParamStorer {
	return singleParamStorer{key: key, value: value, safe: true}
}

// NewUnsafeParamStorer returns a new ParamStorer that stores the provided parameters as UnsafeParams.
func NewUnsafeParamStorer(unsafeParams map[string]interface{}) ParamStorer {
	return NewSafeAndUnsafeParamStorer(nil, unsafeParams)
}

// NewUnsafeParam returns a new ParamStorer that stores a single unsafe parameter.
func NewUnsafeParam(key string, value interface{}) ParamStorer {
	return singleParamStorer{key: key, value: value, safe: false}
}

// NewSafeAndUnsafeParamStorer returns a new ParamStorer that stores the provided safe parameters as SafeParams and the
// unsafe parameters as UnsafeParams. If the safeParams and unsafeParams have any keys in common, the key/value pairs in
// the unsafeParams will be used (the conflicting key/value pairs provided by safeParams will be ignored).
func NewSafeAndUnsafeParamStorer(safeParams, unsafeParams map[string]interface{}) ParamStorer {
	storer := &mapParamStorer{}
	for k, v := range safeParams {
		storer.putSafeParam(k, v)
	}
	for k, v := range unsafeParams {
		storer.putUnsafeParam(k, v)
	}
	return storer
}

type mapParamStorer struct {
	safeParams   map[string]interface{}
	unsafeParams map[string]interface{}
}

func (m *mapParamStorer) SafeParams() map[string]interface{} {
	if m.safeParams == nil {
		return map[string]interface{}{}
	}
	return m.safeParams
}

func (m *mapParamStorer) UnsafeParams() map[string]interface{} {
	if m.unsafeParams == nil {
		return map[string]interface{}{}
	}
	return m.unsafeParams
}

func (m *mapParamStorer) putSafeParam(k string, v interface{}) {
	if m.safeParams == nil {
		m.safeParams = map[string]interface{}{k: v}
	} else {
		m.safeParams[k] = v
	}
	delete(m.unsafeParams, k)
}

func (m *mapParamStorer) putUnsafeParam(k string, v interface{}) {
	if m.unsafeParams == nil {
		m.unsafeParams = map[string]interface{}{k: v}
	} else {
		m.unsafeParams[k] = v
	}
	delete(m.safeParams, k)
}

func (m *mapParamStorer) copyFrom(storer ParamStorer) {
	if storer == nil {
		return
	}
	// If this is one of our types, we can access the values directly and avoid intermediate map allocations.
	switch st := storer.(type) {
	case singleParamStorer:
		if st.safe {
			m.putSafeParam(st.key, st.value)
		} else {
			m.putUnsafeParam(st.key, st.value)
		}
	case *mapParamStorer:
		for k, v := range st.safeParams {
			m.putSafeParam(k, v)
		}
		for k, v := range st.unsafeParams {
			m.putUnsafeParam(k, v)
		}
	default:
		for k, v := range st.SafeParams() {
			m.putSafeParam(k, v)
		}
		for k, v := range st.UnsafeParams() {
			m.putUnsafeParam(k, v)
		}
	}
}

type singleParamStorer struct {
	key   string
	value interface{}
	safe  bool
}

func (s singleParamStorer) SafeParams() map[string]interface{} {
	if !s.safe {
		return map[string]interface{}{}
	}
	return map[string]interface{}{s.key: s.value}
}

func (s singleParamStorer) UnsafeParams() map[string]interface{} {
	if s.safe {
		return map[string]interface{}{}
	}
	return map[string]interface{}{s.key: s.value}
}
