resource "random_id" "suffix" {
  byte_length = 6
}

module "simple_storage_bucket" {
  source  = "harrisonfang/wrappers/google//modules/storage/bucket"
  version = "~> 0.3.0"

  name   = "simple-storage-bucket-${random_id.suffix.hex}"
  region = "us-west1"
}
