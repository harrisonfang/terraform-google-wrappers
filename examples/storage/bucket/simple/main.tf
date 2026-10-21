module "simple_storage_bucket" {
  source  = "harrisonfang/wrappers/google//modules/storage/bucket"
  version = "~> 0.3.0"

  project = var.project

  name   = "${var.project}-simple-storage-bucket"
  region = "us-west1"
}
