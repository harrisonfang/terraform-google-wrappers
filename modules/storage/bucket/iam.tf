resource "google_storage_bucket_iam_member" "bucket_object_viewer" {
  for_each = toset(var.read_only_members)
  bucket   = google_storage_bucket.bucket.name
  role     = "roles/storage.objectViewer"
  member   = each.key
}

resource "google_storage_bucket_iam_member" "bucket_object_creator" {
  for_each = toset(var.write_once_members)
  bucket   = google_storage_bucket.bucket.name
  role     = "roles/storage.objectCreator"
  member   = each.key
}

resource "google_storage_bucket_iam_member" "bucket_object_admin" {
  for_each = toset(var.read_write_members)
  bucket   = google_storage_bucket.bucket.name
  role     = "roles/storage.objectAdmin"
  member   = each.key
}
