package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopySourceExcludesState(t *testing.T) {
	src, dst := t.TempDir(), filepath.Join(t.TempDir(), "copy")
	require.NoError(t, os.Mkdir(filepath.Join(src, ".terraform"), 0700))
	for _, name := range []string{"main.tf", ".terraform.lock.hcl", "terraform.tfstate", "secret.auto.tfvars", ".terraform/main.tf"} {
		require.NoError(t, os.WriteFile(filepath.Join(src, name), []byte("test"), 0600))
	}
	require.NoError(t, os.MkdirAll(filepath.Join(src, "nested/module"), 0700))
	require.NoError(t, os.WriteFile(filepath.Join(src, "nested/module/main.tf"), []byte("test"), 0600))
	require.NoError(t, os.Symlink(filepath.Join(src, "main.tf"), filepath.Join(src, "link.tf")))
	copySource(t, src, dst)
	require.FileExists(t, filepath.Join(dst, "main.tf"))
	require.FileExists(t, filepath.Join(dst, "nested/module/main.tf"))
	require.FileExists(t, filepath.Join(dst, ".terraform.lock.hcl"))
	for _, name := range []string{"terraform.tfstate", "secret.auto.tfvars", ".terraform", "link.tf"} {
		require.NoFileExists(t, filepath.Join(dst, name))
		require.NoDirExists(t, filepath.Join(dst, name))
	}
}
