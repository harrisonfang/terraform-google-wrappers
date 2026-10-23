# Tests

Integration tests use [CFT blueprint-test](https://github.com/GoogleCloudPlatform/cloud-foundation-toolkit/tree/main/infra/blueprint-test), following the upstream Cloud Storage tests, with an existing dedicated test project instead of a project-factory setup.

## Run

Install Go 1.24 or newer and Terraform 1.9 or newer. From this directory:

```sh
go test -short ./...
```

This compiles the harness and runs local checks without creating cloud resources. Integration tests are opt-in and create three empty regional buckets to check project placement, uniform access, Autoclass at the 30-day boundary, lifecycle rules, and an empty post-apply plan.

For an integration run, configure Application Default Credentials for a dedicated non-production project with billing and the Cloud Storage API already enabled. The caller needs bucket create/get/update/delete permissions in that project (for example, Storage Admin on this dedicated project only). No project creation, service-account keys, project IAM changes, or API enablement are performed by the tests.

```sh
TEST_PROJECT_ID=your-dedicated-test-project RUN_GCP_TESTS=1 \
  go test -v -count=1 -timeout=30m -run '^TestStorageBucket$' .
```

These tests can incur charges. Never point them at a production project. Use short-lived credentials. Cloud integration CI is not configured here.

## Isolation and cleanup

Each run uses a random bucket prefix and separate local Terraform state under `tests/.runs/`. Only Terraform source files are copied, not existing state, variable files, or credentials. Cleanup is registered before apply and runs even when apply or assertions fail. Successful runs are removed; failed runs remain for diagnosis and recovery.

A killed process, Go timeout, machine failure, or failed destroy can still leave resources behind. The log prints the recovery workspace. With the same credentials, retry cleanup in its `tests/fixtures/storage_bucket` directory:

```sh
terraform init
terraform destroy
```

The workspace retains non-secret input variables and state. Do not commit it or share state publicly. Remove the run directory only after confirming cleanup. The tests do not write objects or force-delete buckets with unexpected contents.

The fixture uses local source: CFT module-swapper v0.4.20 cannot resolve the nested `modules/storage/bucket` Registry address. No source rewriting or published release is needed.


These tests do not mutate project-wide settings. Separate state does not serialize project-wide changes; no shared-project lock is configured.

## Adding a module test

Add a fixture under `tests/fixtures/<name>` with `project_id` and `name_prefix` inputs, then call `newFixture(t, "<name>")`. Add assertions with `DefineVerify` (including `DefaultVerify` for the no-drift check) and call `Test()`. CFT owns init, apply, verification and teardown; the shared helper only handles opt-in configuration, workspace isolation and recovery. Keep resource-specific setup and assertions in the fixture and its test. The source copy includes `.tf` files and provider locks, not variable files or state.
