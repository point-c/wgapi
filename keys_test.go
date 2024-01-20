package wgapi

import (
	"crypto/rand"
	"errors"
	"github.com/stretchr/testify/require"
	"io"
	"testing"
	"testing/iotest"
)

func TestNewPreshared(t *testing.T) {
	key, err := NewPreshared()
	require.NoError(t, err)
	require.IsType(t, PresharedKey{}, key)
	// Fails to convert to public
	_, err = key.Public()
	require.Error(t, err)
	require.NotEmpty(t, err.Error())
}

func TestNewPrivate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		key, err := NewPrivate()
		require.NoError(t, err)
		require.IsType(t, PrivateKey{}, key)
	})
	t.Run("error", func(t *testing.T) {
		defer func(r io.Reader) { rand.Reader = r }(rand.Reader)
		rand.Reader = iotest.ErrReader(errors.New("test"))
		_, err := NewPrivate()
		require.Error(t, err)
	})
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
	t.Run("error", func(t *testing.T) {
		defer func(r io.Reader) { rand.Reader = r }(rand.Reader)
		rand.Reader = iotest.ErrReader(errors.New("test"))
		_, _, err := NewPrivatePublic()
		require.Error(t, err)
	})
}
