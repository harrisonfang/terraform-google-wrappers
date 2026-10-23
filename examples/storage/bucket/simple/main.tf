module "simple_storage_bucket" {
  source  = "harrisonfang/wrappers/google//modules/storage/bucket"
  version = "~> 0.3.2"

  project = var.project
  name    = "simple-storage-bucket-${random_id.value.hex}"
  region  = "us-west1"
}
