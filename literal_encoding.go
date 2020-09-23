package cache

import (
	"reflect"
	"strconv"
)

var _ Encoding = &literalEncoding{}

// NewLiteralEncoding is an encoding that will try its best to store the data as is,
// but fallback on another encoder if not possible.
func NewLiteralEncoding(fallback Encoding) *literalEncoding {
	return &literalEncoding{fallback: fallback}
}

var translators = map[reflect.Kind]*translator{}

type translator struct {
	encode func(val reflect.Value) []byte
	decode func(data []byte, target reflect.Value) error
}

func init() {
	uintTranslator := &translator{
		func(val reflect.Value) []byte { return []byte(strconv.FormatUint(val.Uint(), 10)) },
		func(data []byte, target reflect.Value) error {
			if val, err := strconv.ParseUint(string(data), 10, 64); err == nil {
				target.SetUint(val)
				return nil
			} else {
				return err
			}
		},
	}
	translators[reflect.Uint] = uintTranslator
	translators[reflect.Uint8] = uintTranslator
	translators[reflect.Uint16] = uintTranslator
	translators[reflect.Uint32] = uintTranslator
	translators[reflect.Uint64] = uintTranslator

	intTranslator := &translator{
		func(val reflect.Value) []byte { return []byte(strconv.FormatInt(val.Int(), 10)) },
		func(data []byte, target reflect.Value) error {
			if val, err := strconv.ParseInt(string(data), 10, 64); err == nil {
				target.SetInt(val)
				return nil
			} else {
				return err
			}
		},
	}
	translators[reflect.Int] = intTranslator
	translators[reflect.Int8] = intTranslator
	translators[reflect.Int16] = intTranslator
	translators[reflect.Int32] = intTranslator
	translators[reflect.Int64] = intTranslator

	floatTranslator := &translator{
		func(val reflect.Value) []byte { return []byte(strconv.FormatFloat(val.Float(), 'f', -1, 64)) },
		func(data []byte, target reflect.Value) error {
			if val, err := strconv.ParseFloat(string(data), 64); err == nil {
				target.SetFloat(val)
				return nil
			} else {
				return err
			}
		},
	}
	translators[reflect.Float32] = floatTranslator
	translators[reflect.Float64] = floatTranslator

	translators[reflect.String] = &translator{
		func(val reflect.Value) []byte { return []byte(val.String()) },
		func(data []byte, target reflect.Value) error {
			target.SetString(string(data))
			return nil
		},
	}

	translators[reflect.Bool] = &translator{
		func(val reflect.Value) []byte { return []byte(strconv.FormatBool(val.Bool())) },
		func(data []byte, target reflect.Value) error {
			if val, err := strconv.ParseBool(string(data)); err == nil {
				target.SetBool(val)
				return nil
			} else {
				return err
			}
		},
	}
}

type literalEncoding struct {
	fallback Encoding
}

func (e *literalEncoding) Encode(data interface{}) ([]byte, error) {
	value := reflect.ValueOf(data)
	if t, ok := translators[value.Kind()]; ok {
		return t.encode(value), nil
	}

	return e.fallback.Encode(data)
}

func (e *literalEncoding) Decode(b []byte, data interface{}) error {
	if !isPointer(data) {
		return ErrNotAPointer
	}

	kind := reflect.ValueOf(data).Elem().Kind()
	if t, ok := translators[kind]; ok {
		return t.decode(b, reflect.ValueOf(data).Elem())
	}

	return e.fallback.Decode(b, data)
}
