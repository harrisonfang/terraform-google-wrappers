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

These tests can incur charges. Never point them at a production project. Use short-lived credentials; a future CI integration should use Workload Identity Federation, a narrowly scoped test identity, and trusted-code approval rather than credentials on untrusted pull requests. Cloud integration CI is deliberately not configured here.

## Isolation and cleanup

Each run uses a random bucket prefix and separate local Terraform state under `tests/.runs/`. Only Terraform source files are copied, not existing state, variable files, or credentials. Cleanup is registered before apply and runs even when apply or assertions fail. Successful runs are removed; failed runs remain for diagnosis and recovery.

A killed process, Go timeout, machine failure, or failed destroy can still leave resources behind. The log prints the recovery workspace. With the same credentials, retry cleanup in its `tests/fixtures/storage_bucket` directory:

```sh
terraform init
terraform destroy
```

The workspace retains non-secret input variables and state. Do not commit it or share state publicly. Remove the run directory only after confirming cleanup. The tests do not write objects or force-delete buckets with unexpected contents.

The fixture references the local module directly and does not depend on a published release or the examples PR. CFT module-swapper v0.4.20 discovers only one level under `modules/`, so it cannot currently resolve this repository's nested `modules/storage/bucket` Registry address. When nested-module support is available, use it on a temporary checkout to test Registry examples without editing their published sources; no custom source rewriting is used here.

## Project-wide tests

The shared project is appropriate for independent bucket resources, not tests that mutate project-wide IAM, APIs, or default networking. Random resource names and separate state do not serialize project-wide changes.

When such tests are added, choose isolation explicitly: either provision a disposable project only for that test, or serialize every actor using the shared project with a common external lock. A GitHub Actions concurrency group only coordinates participating workflows in the same repository; Go parallelism settings and Terraform state locking do not coordinate independent state files or local runs.

An opt-in project-factory fixture should use separate setup state and a separately authorized identity with billing/project-creation permissions. Register setup cleanup before apply, destroy child resources before the project, use a random project ID, and explicitly allow deletion (`deletion_policy = "DELETE"`; the factory defaults to `PREVENT`). Missing factory configuration must not silently fall back to the shared project. A project-wrapper test should test its own project creation rather than factory-creating the same object. No disposable-project provisioner or global lock is implemented until a test needs one.
