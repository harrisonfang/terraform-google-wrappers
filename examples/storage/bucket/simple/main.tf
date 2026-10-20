module "simple_storage_bucket" {
  source  = "harrisonfang/wrappers/google//modules/storage/bucket"
  version = "~> 0.3.0"

  name   = "${var.project_id}-simple-storage-bucket"
  region = "us-west1"
}
