package form

import (
	"fmt"
	"reflect"
)

// InvalidFieldError indicates the data for a field was invalid.
type InvalidFieldError struct {
	Field string
	Err   error
}

// NewInvalidFieldError returns a InvalidFieldError for the field f and error err.
func NewInvalidFieldError(f string, err error) *InvalidFieldError {
	return &InvalidFieldError{Field: f, Err: err}
}

func (e *InvalidFieldError) Error() string {
	return fmt.Sprintf("form: invalid field %q err: %v", e.Field, e.Err)
}

func (e *InvalidFieldError) PublicError() string {
	return fmt.Sprintf("form: invalid field %q", e.Field)
}

// InvalidParameterError indicates a parameter was passed but of an invalid type
type InvalidParameterError struct {
	Parameter string
	Type      reflect.Type
	Required  string
	source    string
}

// NewInvalidFormParameterError creates the associated named error type implicitly with form source and populates it.
func NewInvalidFormParameterError(p string, t reflect.Type, r string) *InvalidParameterError {
	return &InvalidParameterError{Parameter: p, Type: t, Required: r, source: "form"}
}

// NewInvalidJSONParameterError creates the associated named error type implicitly with form source and populates it.
func NewInvalidJSONParameterError(p string, t reflect.Type, r string) *InvalidParameterError {
	return &InvalidParameterError{Parameter: p, Type: t, Required: r, source: "json"}
}

func (e *InvalidParameterError) Error() string {
	return fmt.Sprintf("%v: parameter %q must be %v, not %v", e.source, e.Parameter, e.Required, e.Type.Kind())
}

// MissingFieldError indicates a field was missing.
type MissingFieldError struct {
	Field  string
	source string
}

// NewMissingFormFieldError returns a MissingFieldError for the field f with form source.
func NewMissingFormFieldError(f string) *MissingFieldError {
	return &MissingFieldError{Field: f, source: "form"}
}

// NewMissingJSONFieldError returns a MissingFieldError for the field f with json source.
func NewMissingJSONFieldError(f string) *MissingFieldError {
	return &MissingFieldError{Field: f, source: "json"}
}

func (e *MissingFieldError) Error() string {
	return fmt.Sprintf("%v: missing field %q", e.source, e.Field)
}

func (e *MissingFieldError) PublicError() string {
	return fmt.Sprintf("%v: missing field %q", e.source, e.Field)
}

// UnsupportedTypeError indicates a struct field had a type we couldn't handle.
type UnsupportedTypeError struct {
	Type   reflect.Type
	Source reflect.Type
	Field  string
	src    string
}

// NewUnsupportedFormTypeError returns a new UnsupportedTypeError for type t, source s and field f.
func NewUnsupportedFormTypeError(t reflect.Type, s reflect.Type, f string) *UnsupportedTypeError {
	return &UnsupportedTypeError{Type: t, Source: s, Field: f, src: "form"}
}

// NewUnsupportedJSONTypeError returns a new UnsupportedTypeError for type t, source s and field f.
func NewUnsupportedJSONTypeError(t reflect.Type, s reflect.Type, f string) *UnsupportedTypeError {
	return &UnsupportedTypeError{Type: t, Source: s, Field: f, src: "json"}
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("%v: unsupported type %q for field %v.%v", e.src, e.Type.Name(), e.Source.String(), e.Field)
}
