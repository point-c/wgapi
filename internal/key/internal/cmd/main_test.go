package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_main(t *testing.T) {
	defer func(args ...string) { os.Args = args }(os.Args...)
	td := t.TempDir()
	gen := filepath.Join(td, "keys_generated.go")
	genTest := filepath.Join(td, "keys_generated_test.go")
	cfg := filepath.Join(td, "keys.yml")
	require.NoError(t, os.WriteFile(cfg, []byte("foo: bar\n"), os.ModePerm))
	os.Args = []string{"test", "-out", gen, "-tout", genTest, "-config", cfg}
	main()
	require.FileExists(t, gen)
	require.FileExists(t, genTest)
}
