package athena_abi

import (
	"math/big"
	"testing"
)

// Mock StarknetType to use in tests.
type mockType struct {
	id string
}

func (m mockType) idStr() string {
	return m.id
}

/*tests for athena_abi */

// TestIntFromString tests the intFromString function for StarknetCoreType.
func TestIntFromString(t *testing.T) {
	tests := []struct {
		input    string
		expected StarknetCoreType
	}{
		{"u8", U8},
		{"u16", U16},
		{"u32", U32},
		{"u64", U64},
		{"u128", U128},
		{"u256", U256},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := intFromString(tt.input)
			if err != nil {
				t.Fatalf("error converting string to StarknetCoreType: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestMaxValue tests the maxValue method for StarknetCoreType.
func TestMaxValue(t *testing.T) {
	tests := []struct {
		input       StarknetCoreType
		expectedStr string
	}{
		{U8, "255"},
		{U16, "65535"},
		{U32, "4294967295"},
		{U256, "115792089237316195423570985008687907853269984665640564039457584007913129639935"},
		{EthAddress, "1461501637330902918203684832716283019655932542975"},
	}

	for _, tt := range tests {
		t.Run(tt.input.String(), func(t *testing.T) {
			maxValue, err := tt.input.maxValue()
			if err != nil {
				t.Fatalf("error while getting maxValue: %v", err)
			}
			expected := new(big.Int)
			expected.SetString(tt.expectedStr, 10)

			if maxValue.Cmp(expected) != 0 {
				t.Errorf("expected %s, got %s", expected.String(), maxValue.String())
			}
		})
	}
}

func TestStarknetArray(t *testing.T) {
	array := StarknetArray{
		InnerType: mockType{id: "Felt"},
	}
	expected := "[Felt]"
	if result := array.idStr(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestStarknetOption(t *testing.T) {
	option := StarknetOption{
		InnerType: mockType{id: "U8"},
	}
	expected := "Option[U8]"
	if result := option.idStr(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestStarknetNonZero(t *testing.T) {
	nonZero := StarknetNonZero{
		InnerType: mockType{id: "U32"},
	}
	expected := "NonZero[U32]"
	if result := nonZero.idStr(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestStarknetEnum(t *testing.T) {
	enum := StarknetEnum{
		Name: "MyEnum",
		Variants: []struct {
			Name string
			Type StarknetType
		}{
			{"Variant1", mockType{id: "U16"}},
			{"Variant2", mockType{id: "NoneType"}},
		},
	}
	expected := "Enum[Variant1:U16,Variant2]"
	if result := enum.idStr(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestStarknetTuple(t *testing.T) {
	tuple := StarknetTuple{
		Members: []StarknetType{
			mockType{id: "U16"},
			mockType{id: "U32"},
		},
	}
	expected := "(U16,U32)"
	if result := tuple.idStr(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAbiParameter(t *testing.T) {
	param := AbiParameter{
		Name: "param1",
		Type: mockType{id: "Felt"},
	}
	expected := "param1:Felt"
	if result := param.idStr(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestStarknetStruct(t *testing.T) {
	starknetStruct := StarknetStruct{
		Name: "MyStruct",
		Members: []AbiParameter{
			{"field1", mockType{id: "U8"}},
			{"field2", mockType{id: "Felt"}},
		},
	}
	expected := "{field1:U8,field2:Felt}"
	if result := starknetStruct.idStr(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

// Testing cases of unknown and invalid type values //

// Unknown StarknetCoreType
func TestZeroValueStarknetType(t *testing.T) {
	// Zero or invalid StarknetCoreType value
	var zeroValue StarknetCoreType

	// Testing idStr for a zero value
	expectedIDStr := "Unknown" // Since zero doesn't match any valid enum
	if result := zeroValue.idStr(); result != expectedIDStr {
		t.Errorf("expected %s, got %s", expectedIDStr, result)
	}

	// Testing maxValue for a zero value
	_, err := zeroValue.maxValue()
	if err == nil {
		t.Error("expected error for zero-value maxValue but got nil")
	}
}

// Inavlid types
func TestIntFromStringInvalid(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"invalid", "invalid integer type: invalid"},
		{"", "invalid integer type: "},
		{"123", "invalid integer type: 123"},
		{"unknown", "invalid integer type: unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := intFromString(tt.input)
			if err == nil || err.Error() != tt.expected {
				t.Errorf("expected error %s, got %v", tt.expected, err)
			}
		})
	}
}

func TestMaxValueInvalidType(t *testing.T) {
	invalidType := StarknetCoreType(999) // N.A value

	maxValue, err := invalidType.maxValue()
	if err == nil {
		t.Error("expected error for invalid type, but got nil")
	}

	if maxValue != nil {
		t.Errorf("expected nil maxValue for invalid type, got %v", maxValue)
	}
}
