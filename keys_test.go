package wgapi

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPreshared(t *testing.T) {
	key, err := NewPreshared()
	require.NoError(t, err)
	require.IsType(t, PresharedKey{}, key)

	_, err = key.Public()
	require.Error(t, err)
	require.NotEmpty(t, err.Error())
}

func TestNewPrivate(t *testing.T) {
	key, err := NewPrivate()
	require.NoError(t, err)
	require.IsType(t, PrivateKey{}, key)
}

func TestNewPrivatePublic(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		pr, pu, err := NewPrivatePublic()
		require.NoError(t, err)
		require.IsType(t, PublicKey{}, pu)
		require.IsType(t, PrivateKey{}, pr)

		pu, err = pu.Public()
		require.NoError(t, err)
		require.IsType(t, PublicKey{}, pu)

		b, err := pu.MarshalText()
		require.NoError(t, err)
		require.NotEmpty(t, b)

		err = pu.UnmarshalText(b)
		require.NoError(t, err)

		err = pu.UnmarshalText(b[:1])
		require.Error(t, err)
	})
}
