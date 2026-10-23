package tests

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageBucket(t *testing.T) {
	if testing.Short() || os.Getenv("RUN_GCP_TESTS") != "1" {
		t.Skip("set RUN_GCP_TESTS=1 and TEST_PROJECT_ID to run billable integration tests")
	}
	project := os.Getenv("TEST_PROJECT_ID")
	require.NotEmpty(t, project, "TEST_PROJECT_ID must name a dedicated, existing test project")

	// Keep failed runs on disk: their state is needed to retry a failed destroy.
	require.NoError(t, os.MkdirAll(".runs", 0700))
	run, err := os.MkdirTemp(".runs", "storage-bucket-")
	require.NoError(t, err)
	run, err = filepath.Abs(run)
	require.NoError(t, err)
	t.Logf("Recovery workspace: %s", run)
	copyTree(t, "../modules/storage/bucket", filepath.Join(run, "modules/storage/bucket"))
	fixture := filepath.Join(run, "tests/fixtures/storage_bucket")
	copyTree(t, "fixtures/storage_bucket", fixture)
	// Prevent inherited Terraform settings from sharing state/data across runs.
	t.Setenv("TF_DATA_DIR", filepath.Join(fixture, ".terraform"))
	t.Setenv("TF_WORKSPACE", "default")

	id := make([]byte, 8)
	_, err = rand.Read(id)
	require.NoError(t, err)
	prefix := "tgw-test-" + hex.EncodeToString(id)
	vars := map[string]interface{}{"project_id": project, "name_prefix": prefix}
	// Save only non-secret inputs so cleanup never needs to reconstruct arguments.
	writeVars(t, fixture, project, prefix)
	bp := tft.NewTFBlueprintTest(t, tft.WithTFDir(fixture), tft.WithVars(vars))
	initialized := false
	t.Cleanup(func() {
		if initialized {
			if _, err := terraform.DestroyE(t, bp.GetTFOptions()); err != nil {
				t.Errorf("destroy failed; retain %s and retry terraform destroy there: %v", fixture, err)
				return
			}
		}
		if !t.Failed() {
			assert.NoError(t, os.RemoveAll(run))
		}
	})
	a := assert.New(t)
	bp.Init(a)
	initialized = true
	bp.Apply(a)
	bp.DefaultVerify(a)
	buckets := bp.GetJsonOutput("buckets")
	for _, key := range []string{"default", "short", "long"} {
		bucket := buckets.Get(key)
		a.Equal(project, bucket.Get("project").String())
		a.Equal(prefix+"-"+key, bucket.Get("name").String())
		a.Equal("US-CENTRAL1", bucket.Get("location").String())
		a.True(bucket.Get("uniform_bucket_level_access").Bool())
		a.Equal(key != "short", bucket.Get("autoclass.0.enabled").Bool())
		if key != "short" {
			a.Equal("ARCHIVE", bucket.Get("autoclass.0.terminal_storage_class").String())
		}
		if key == "default" {
			a.Empty(bucket.Get("lifecycle_rule").Array())
		} else {
			age := int64(30)
			if key == "long" {
				age = 31
			}
			a.Equal(age, bucket.Get("lifecycle_rule.0.condition.0.age").Int())
			a.Equal("Delete", bucket.Get("lifecycle_rule.0.action.0.type").String())
		}
	}
}
