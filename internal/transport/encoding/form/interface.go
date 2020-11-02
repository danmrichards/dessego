package form

// Unmarshaler is the interface implemented by objects that can unmarshal themselves into valid object from a form field.
type Unmarshaler interface {
	UnmarshalForm(string) error
}

// Empty is the interface implemented by objects which can determine if the
// struct is empty
type Empty interface {
	IsEmpty() bool
}
