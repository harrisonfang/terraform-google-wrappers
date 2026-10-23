package tests

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newFixture prepares any fixture using the shared project_id/name_prefix inputs.
// The caller adds resource-specific verification and invokes CFT's Test lifecycle.
func newFixture(t *testing.T, name string) *tft.TFBlueprintTest {
	t.Helper()
	if testing.Short() || os.Getenv("RUN_GCP_TESTS") != "1" {
		t.Skip("set RUN_GCP_TESTS=1 and TEST_PROJECT_ID to run billable integration tests")
	}
	project := os.Getenv("TEST_PROJECT_ID")
	require.NotEmpty(t, project, "TEST_PROJECT_ID must name a dedicated, existing test project")
	require.NoError(t, os.MkdirAll(".runs", 0700))
	run, err := os.MkdirTemp(".runs", name+"-")
	require.NoError(t, err)
	run, err = filepath.Abs(run)
	require.NoError(t, err)
	t.Logf("Recovery workspace: %s", run)
	copySource(t, "../modules", filepath.Join(run, "modules"))
	fixture := filepath.Join(run, "tests/fixtures", name)
	copySource(t, filepath.Join("fixtures", name), fixture)
	t.Setenv("TF_DATA_DIR", filepath.Join(fixture, ".terraform"))
	t.Setenv("TF_WORKSPACE", "default")
	id := make([]byte, 8)
	_, err = rand.Read(id)
	require.NoError(t, err)
	vars := map[string]interface{}{"project_id": project, "name_prefix": "tgw-test-" + hex.EncodeToString(id)}
	// Persist non-secret inputs with the recovery state.
	data, err := json.Marshal(vars)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(fixture, "test.auto.tfvars.json"), data, 0600))
	bp := tft.NewTFBlueprintTest(t, tft.WithTFDir(fixture), tft.WithVars(vars))
	bp.DefineTeardown(func(a *assert.Assertions) {
		// CFT defers this before apply. A failed destroy stops here, preserving state.
		bp.DefaultTeardown(a)
		if !t.Failed() {
			require.NoError(t, os.RemoveAll(run))
		}
	})
	return bp
}

// Use upstream copying, with an allowlist that excludes state, variables and links.
func copySource(t *testing.T, src, dst string) {
	t.Helper()
	require.NoError(t, os.MkdirAll(dst, 0700))
	require.NoError(t, files.CopyFolderContentsWithFilter(src, dst, func(path string) bool {
		info, err := os.Lstat(path)
		require.NoError(t, err)
		if info.Mode()&os.ModeSymlink != 0 || info.Name() == ".terraform" {
			return false
		}
		return info.IsDir() || filepath.Ext(path) == ".tf" || info.Name() == ".terraform.lock.hcl"
	}))
}
