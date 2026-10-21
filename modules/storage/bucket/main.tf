resource "google_storage_bucket" "bucket" {
  name     = var.name
  location = var.region

  uniform_bucket_level_access = true

  autoclass {
    enabled                = var.object_ttl_days == null ? true : var.object_ttl_days > 30
    terminal_storage_class = "ARCHIVE"
  }

  dynamic "lifecycle_rule" {
    for_each = var.object_ttl_days == null ? [] : [1]
    content {
      condition {
        age = var.object_ttl_days
      }
      action {
        type = "Delete"
      }
    }
  }
}
