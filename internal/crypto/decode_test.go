package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestDecrypter_Decrypt(t *testing.T) {
	d, err := NewDecrypter(DefaultAESKey)
	if err != nil {
		t.Fatal(err)
	}

	plain := []byte("hello world")
	enc, err := testEnc(plain)
	if err != nil {
		t.Fatal(err)
	}

	if dec := d.Decrypt(enc); !reflect.DeepEqual(plain, dec) {
		t.Fatalf("expected %q got %q", plain, dec)
	}
}

func testEnc(plain []byte) ([]byte, error) {
	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	bl := len(plain)
	for bl%aes.BlockSize != 0 {
		bl++
	}

	// Padding delta.
	d := bl - len(plain)
	for i := 0; i < d; i++ {
		plain = append(plain, byte(d))
	}

	block, err := aes.NewCipher([]byte(DefaultAESKey))
	if err != nil {
		return nil, fmt.Errorf("new cipher: %w", err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	enc := make([]byte, aes.BlockSize+len(plain))
	iv := enc[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("read iv: %w", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(enc[aes.BlockSize:], plain)

	return enc, nil
}
