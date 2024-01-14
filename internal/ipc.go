package internal

import (
	"fmt"
)

type KeyValue interface {
	Key() string
	fmt.Stringer
}

type TestKey struct{}

const TestKeyString = "test"

func (TestKey) String() string { return TestKeyString }
