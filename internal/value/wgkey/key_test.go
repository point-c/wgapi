package wgkey

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeytype(t *testing.T) {
	Public{}.keyType()
	Private{}.keyType()
	PreShared{}.keyType()
}

func TestErrInvalidKeyType(t *testing.T) {
	testErrType[Public](t)
	testErrType[Private](t)
	testErrType[PreShared](t)
}

func testErrType[T Type](t *testing.T) {
	t.Helper()
	var errPrivate ErrInvalidKeyType[T]
	require.ErrorContains(t, &errPrivate, fmt.Sprintf("%T", *new(T)))
}
