package tests

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// Copy source only, never developer state, credentials, or provider caches.
func copyTree(t *testing.T, src, dst string) {
	t.Helper()
	require.NoError(t, filepath.WalkDir(src, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			if entry.Name() == ".terraform" {
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(dst, rel), 0700)
		}
		if entry.Type()&os.ModeSymlink != 0 || (filepath.Ext(path) != ".tf" && entry.Name() != ".terraform.lock.hcl") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), data, 0600)
	}))
}

func writeVars(t *testing.T, fixture, project, prefix string) {
	t.Helper()
	data, err := json.Marshal(map[string]string{"project_id": project, "name_prefix": prefix})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(fixture, "test.auto.tfvars.json"), data, 0600))
}

func TestCopyTreeExcludesState(t *testing.T) {
	src, dst := t.TempDir(), filepath.Join(t.TempDir(), "copy")
	require.NoError(t, os.Mkdir(filepath.Join(src, ".terraform"), 0700))
	for _, name := range []string{"main.tf", ".terraform.lock.hcl", "terraform.tfstate", "secret.auto.tfvars", ".terraform/main.tf"} {
		require.NoError(t, os.WriteFile(filepath.Join(src, name), []byte("test"), 0600))
	}
	copyTree(t, src, dst)
	require.FileExists(t, filepath.Join(dst, "main.tf"))
	require.FileExists(t, filepath.Join(dst, ".terraform.lock.hcl"))
	for _, name := range []string{"terraform.tfstate", "secret.auto.tfvars", ".terraform"} {
		require.NoFileExists(t, filepath.Join(dst, name))
		require.NoDirExists(t, filepath.Join(dst, name))
	}
}
