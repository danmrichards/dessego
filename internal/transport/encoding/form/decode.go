package form

import (
	"bytes"
	"database/sql"
	"encoding"
	"encoding/json"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// Decoder decodes data from a url.Values
type Decoder struct {
	v url.Values
	// OmitEmpty = true will omit all fields which have the empty value for their type regardless of any tags
	OmitEmpty bool
	// LowerKeys = true will ensure all value keys are lowercase regardless of any tags
	LowerKeys bool
	// TagName is the name of the tag to use if non empty, otherwise DefaultTagName is used.
	TagName string
	err     error
}

// NewDecoder returns a new decoder that writes to v.
func NewDecoder(v url.Values) *Decoder {
	return &Decoder{v: v, TagName: DefaultTagName}
}

// Decode writes the values decoded from url.Values to a struct or struct pointer
//
// The only supported types for fields are:
// * int, int8, int16, int32, int64
// * uint, uint8, uint16, uint32, uint64
// * float32, float64
// * string
// * bool
// * struct
// * Ptr to one of the above
//
// Each exported field becomes a key of the Values unless:
//   - the field's tag name is "-", or
//   - the field is empty and its tag specifies the "omitempty" option or the decoder has OmitEmpty set to true
//   - the field is a Ptr which is nil
//
// Each exported field key will be the name of the field unless:
//   - the field's tag name is not "-"
//   - the field's tag specifies the "lowerkey" option or the decoder has LowerKeys set to true
//
// The supported tag flags are:
// * omitempty - value is untouched if the form didn't contain a corresponding element.
// * required - an error is returned if the form didn't contain the corresponding element or if its value was empty string.
// * json - uses json.Unmarshal to decode the value, only valid for struct ptrs.
//
// The empty values are false, 0, any nil pointer or interface value, and any array, slice, map, or string of length zero.
func (dec *Decoder) Decode(v interface{}) error {
	if v == nil {
		return nil
	}

	if dec.TagName == "" {
		dec.TagName = DefaultTagName
	}

	val := reflect.ValueOf(v)
	var invalid bool
	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
		if val.Kind() != reflect.Struct {
			invalid = true
		}
	case reflect.Struct:
		if !val.CanAddr() {
			invalid = true
		}
		// OK
	default:
		invalid = true
	}

	if invalid {
		dec.err = NewInvalidFormParameterError("v", val.Type(), "a struct or pointer to a struct")
		return dec.err
	}

	return dec.decodeStruct(val)
}

func (dec *Decoder) decodeStruct(val reflect.Value) error {
	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if !f.CanInterface() {
			// Skip unexported fields
			continue
		}

		s := t.Field(i)

		tag := s.Tag.Get(dec.TagName)
		if tag == SkipField {
			continue
		}
		name, opts := ParseTag(tag)
		if name == SkipField {
			continue
		}

		k := s.Name
		if name != "" {
			k = name
		}
		if dec.LowerKeys || opts.Contains("lowerkey") {
			k = strings.ToLower(k)
		}

		v := dec.v.Get(k)
		if v == "" {
			if opts.Contains("required") {
				return NewMissingFormFieldError(k)
			}
		}

		switch f.Kind() {
		case reflect.Ptr, reflect.Interface:
			if v == "" {
				// Set to nil
				f.Set(reflect.Zero(f.Type()))
				continue
			} else if f.IsNil() {
				// Create an element for us to set
				f.Set(reflect.New(f.Type().Elem()))
			}
			f = f.Elem()
		default:
			if v == "" && !s.Anonymous {
				// No value so no action needed, we leave any value as it was.
				continue
			}
		}

		unsupported := false
		switch f.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.ParseInt(v, 10, f.Type().Bits())
			if err != nil {
				return NewInvalidFieldError(k, err)
			}
			if dec.skipEmpty(opts, i) {
				continue
			}

			f.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			i, err := strconv.ParseUint(v, 10, f.Type().Bits())
			if err != nil {
				return NewInvalidFieldError(k, err)
			}
			if dec.skipEmpty(opts, i) {
				continue
			}

			f.SetUint(i)
		case reflect.Float32, reflect.Float64:
			fv, err := strconv.ParseFloat(v, f.Type().Bits())
			if err != nil {
				return NewInvalidFieldError(k, err)
			}
			if dec.skipEmpty(opts, fv) {
				continue
			}

			f.SetFloat(fv)
		case reflect.String:
			if dec.skipEmpty(opts, v) {
				continue
			}

			f.SetString(v)
		case reflect.Bool:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return NewInvalidFieldError(k, err)
			}
			if dec.skipEmpty(opts, b) {
				continue
			}

			f.SetBool(b)
		case reflect.Array, reflect.Slice, reflect.Map:
			switch {
			case opts.Contains("json"):
				i := f.Addr().Interface()
				if err := json.Unmarshal([]byte(v), i); err != nil {
					return NewInvalidFieldError(k, err)
				}
				f.Set(reflect.ValueOf(i).Elem())
			case f.Kind() == reflect.Slice:
				et := f.Type().Elem()
				kind := et.Kind()
				switch kind {
				case reflect.String:
					// Slice of Strings
					s := len(dec.v[k])
					slice := reflect.MakeSlice(reflect.SliceOf(et), s, s)
					for i, val := range dec.v[k] {
						slice.Index(i).SetString(val)
					}
					f.Set(slice)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					// Slice of Ints
					s := len(dec.v[k])
					b := et.Bits()
					slice := reflect.MakeSlice(reflect.SliceOf(et), s, s)
					for i, val := range dec.v[k] {
						intVal, err := strconv.ParseInt(val, 10, b)
						if err != nil {
							return NewInvalidFieldError(k, err)
						}
						slice.Index(i).SetInt(intVal)
					}
					f.Set(slice)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					// Slice of Uints
					s := len(dec.v[k])
					b := et.Bits()
					slice := reflect.MakeSlice(reflect.SliceOf(et), s, s)
					for i, val := range dec.v[k] {
						intVal, err := strconv.ParseUint(val, 10, b)
						if err != nil {
							return NewInvalidFieldError(k, err)
						}
						slice.Index(i).SetUint(intVal)
					}
					f.Set(slice)
				default:
					// TODO: Implement full support
					unsupported = true
				}
			default:
				// TODO: Implement full support
				unsupported = true
			}
		case reflect.Struct:
			i := f.Addr().Interface()
			if opts.Contains("json") {
				if err := json.Unmarshal([]byte(v), i); err != nil {
					return NewInvalidFieldError(k, err)
				}
				f.Set(reflect.ValueOf(i).Elem())
			} else {
				switch i := i.(type) {
				case Unmarshaler:
					var err error
					if err = i.UnmarshalForm(v); err != nil {
						return NewInvalidFieldError(k, err)
					}
					f.Set(reflect.ValueOf(i).Elem())
				case sql.Scanner:
					if err := i.Scan(v); err != nil {
						return NewInvalidFieldError(k, err)
					}
				case encoding.TextUnmarshaler:
					var buf bytes.Buffer

					buf.WriteString(v)
					if err := i.UnmarshalText(buf.Bytes()); err != nil {
						return NewInvalidFieldError(k, err)
					}
				default:
					if err := dec.decodeStruct(f); err != nil {
						return NewInvalidFieldError(k, err)
					}
				}
			}
		default:
			unsupported = true
		}

		if unsupported {
			dec.err = NewUnsupportedFormTypeError(f.Type(), t, s.Name)
			return dec.err
		}
	}

	return nil
}

func (dec *Decoder) skipEmpty(opts TagOptions, i interface{}) bool {
	if (dec.OmitEmpty || opts.Contains("omitempty")) && isEmptyValue(reflect.ValueOf(i)) {
		return true
	}

	return false
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		if v.CanAddr() {
			v = v.Addr()
		}
		if v.CanInterface() {
			if i, ok := v.Interface().(Empty); ok {
				return i.IsEmpty()
			}
		}
	}
	return false
}
