package form

import (
	"net/url"
	"reflect"
	"testing"
)

type testDecodeStruct struct {
	StructStr string
	StructInt int
}

type decodeResult struct {
	Int   int
	Int8  int8
	Int16 int16
	Int32 int32
	Int64 int64

	Uint   uint
	Uint8  uint8
	Uint16 uint16
	Uint32 uint32
	Uint64 uint64

	Float32 float32
	Float64 float64

	IntSlice    []int
	ByteSlice   []byte
	Uint8Slice  []uint8
	StringSlice []string

	String string
	Bool   bool

	PtrInt      *int
	PtrIntUnset *int

	StructPtr *testDecodeStruct `form:",json"`
}

func TestDecoder_Decode(t *testing.T) {
	src := url.Values{
		"Int":   []string{"1"},
		"Int8":  []string{"2"},
		"Int16": []string{"3"},
		"Int32": []string{"4"},
		"Int64": []string{"5"},

		"Uint":   []string{"6"},
		"Uint8":  []string{"7"},
		"Uint16": []string{"8"},
		"Uint32": []string{"9"},
		"Uint64": []string{"10"},

		"Float32": []string{"1.1"},
		"Float64": []string{"1.2"},

		"IntSlice":    []string{"0", "1"},
		"ByteSlice":   []string{"1", "2"},
		"Uint8Slice":  []string{"3", "4"},
		"StringSlice": []string{"a", "b"},

		"String": []string{"test"},
		"Bool":   []string{"true"},

		"PtrInt": []string{"1"},

		"StructPtr": []string{`{"StructStr":"test","StructInt":2}`},
	}

	exp := decodeResult{
		Int:   1,
		Int8:  2,
		Int16: 3,
		Int32: 4,
		Int64: 5,

		Uint:   6,
		Uint8:  7,
		Uint16: 8,
		Uint32: 9,
		Uint64: 10,

		Float32: 1.1,
		Float64: 1.2,

		IntSlice:    []int{0, 1},
		ByteSlice:   []byte{1, 2},
		Uint8Slice:  []uint8{3, 4},
		StringSlice: []string{"a", "b"},

		String: "test",
		Bool:   true,

		PtrInt:      &[]int{1}[0],
		PtrIntUnset: nil,

		StructPtr: &testDecodeStruct{StructStr: "test", StructInt: 2},
	}

	var d decodeResult
	if err := NewDecoder(src).Decode(&d); err != nil {
		t.Errorf("decode: %v", err)
		return
	}

	// Compare the values pointed to not the value of the pointer
	if (exp.PtrInt != nil && d.PtrInt == nil) || (exp.PtrInt == nil && d.PtrInt != nil) {
		t.Errorf("expected PtrInt %v, got %v", exp.PtrInt, d.PtrInt)
		return
	} else if *exp.PtrInt != *d.PtrInt {
		t.Errorf("expected PtrInt %v, got %v", *exp.PtrInt, *d.PtrInt)
		return
	}
	d.PtrInt = exp.PtrInt

	// struct containing byte arrays can not be compared, so instead
	// we use deep equals
	if !reflect.DeepEqual(exp, d) {
		t.Errorf("expected %v, got %v", exp, d)
	}
}
