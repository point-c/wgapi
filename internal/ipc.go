package internal

import (
	"fmt"
)

// KeyValue is something represented by a key that can return its value as a string.
type KeyValue interface {
	Key() string
	fmt.Stringer
}

type TestKey struct{}

const TestKeyString = "test"

func (TestKey) String() string { return TestKeyString }
