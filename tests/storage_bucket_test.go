package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageBucket(t *testing.T) {
	bp := newFixture(t, "storage_bucket")
	bp.DefineVerify(func(a *assert.Assertions) {
		bp.DefaultVerify(a)
		vars := bp.GetTFOptions().Vars
		project, prefix := vars["project_id"], vars["name_prefix"].(string)
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
	})
	bp.Test()
}
