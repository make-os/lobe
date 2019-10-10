package types

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// FieldError is used to describe an error concerning an objects field/property
func FieldError(field, err string) error {
	return fmt.Errorf(fmt.Sprintf("field:%s, error:%s", field, err))
}

// FieldErrorWithIndex is used to describe an error concerning an field/property
// of an object contained in list (array or slice).
// If index is -1, it will revert to FieldError
func FieldErrorWithIndex(index int, field, err string) error {
	if index == -1 {
		return FieldError(field, err)
	}
	var fieldArg = "field:%s, "
	if field == "" {
		fieldArg = "%s"
	}
	return fmt.Errorf(fmt.Sprintf("index:%d, "+fieldArg+"error:%s", index, field, err))
}

// ErrStaleSecretRound returns an error about `secretRound` field
// of a tx when the field is not greater than the previous secret round
var ErrStaleSecretRound = func(index int) error {
	return FieldErrorWithIndex(index, "secretRound",
		"must be greater than the previous round")
}

// IsStaleSecretRoundErr checks whether an error is a ErrStaleSecretRound error
func IsStaleSecretRoundErr(err error) bool {
	return strings.Index(err.Error(), "error:must be greater than the previous round") != -1
}

// ErrEarlySecretRound returns an error about `secretRound` field
// of a tx when the field is lower that the expected round.
var ErrEarlySecretRound = func(index int) error {
	return FieldErrorWithIndex(index, "secretRound", "round was generated too early")
}

// IsEarlySecretRoundErr checks whether an error is a ErrEarlySecretRound error
func IsEarlySecretRoundErr(err error) bool {
	return strings.Index(err.Error(), "error:round was generated too early") != -1
}

// HexBytes contains bytes that are encodeable to hex
type HexBytes []byte

// String returns hex string
func (h *HexBytes) String() string {
	return hex.EncodeToString(*h)
}

// HexBytesFromHex returns HexBytes from hex string.
// Panics if hexStr could not be decoded to hex
func HexBytesFromHex(hexStr string) HexBytes {
	bz, err := hex.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}
	return HexBytes(bz)
}
