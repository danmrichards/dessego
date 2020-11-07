package transport

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/danmrichards/dessego/internal/transport/encoding/form"
)

// RequestDecrypter is the interface that wraps the basic Decrypt method.
//
// Decrypt returns a byte slice containing the decrypted contents of the given
// input byte slice.
type RequestDecrypter interface {
	Decrypt([]byte) []byte
}

// DecodeRequest decodes the bytes from a request body into the target, v.
func DecodeRequest(rd RequestDecrypter, data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return InvalidDecodeTargetError{Type: reflect.TypeOf(v)}
	}

	// Demon's Souls sends it's request body as an AES encrypted version of
	// a standard HTTP form POST.
	req := rd.Decrypt(data)

	// Can now use the body as a normal HTTP form.
	vals, err := url.ParseQuery(string(req))
	if err != nil {
		return fmt.Errorf("parse request vals: %w", err)
	}

	if err = form.NewDecoder(vals).Decode(v); err != nil {
		return fmt.Errorf("decode request: %w", err)
	}

	return nil
}
