package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_KeyValue(t *testing.T) {
	require.Equal(t, TestKeyString, TestKey{}.String())
}
