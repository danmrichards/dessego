package transport

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
)

// WriteResponse writes a gamestate response command and data to the given writer.
func WriteResponse(w io.Writer, rt ResponseType, data []byte) error {
	rb, err := buildResponse(rt, data)
	if err != nil {
		return fmt.Errorf("build response: %w", err)
	}

	res := make([]byte, base64.StdEncoding.EncodedLen(len(rb)))
	base64.StdEncoding.Encode(res, rb)

	if _, err := w.Write(res); err != nil {
		return err
	}

	// Game responses require a trailing newline character.
	if _, err := w.Write([]byte("\n")); err != nil {
		return err
	}

	return nil
}

// buildResponse returns a byte slice representing a gamestate server response.
//
// Responses are in the format <CMD_FLAG><DATA_LENGTH><DATA>
func buildResponse(rt ResponseType, data []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Command flag.
	buf.WriteString(fmt.Sprintf("%c", rune(rt)))

	// Data length.
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(data)+5)); err != nil {
		return nil, err
	}

	// Data.
	buf.Write(data)

	return buf.Bytes(), nil
}
