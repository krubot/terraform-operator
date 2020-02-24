
# ---------------------------------------------------------------------------------------------------------------------
# IAM BINDINGS FOR THE CREATED BUCKET
# ---------------------------------------------------------------------------------------------------------------------

resource "google_storage_bucket_iam_binding" "admin" {
  count = length(var.admins) > 0 ? 1 : 0

  bucket  = google_storage_bucket.bucket.name
  role    = "roles/storage.objectAdmin"
  members = var.admins
}

resource "google_storage_bucket_iam_binding" "creator" {
  count = length(var.creators) > 0 ? 1 : 0

  bucket  = google_storage_bucket.bucket.name
  role    = "roles/storage.objectCreator"
  members = var.creators
}

resource "google_storage_bucket_iam_binding" "viewer" {
  count = length(var.viewers) > 0 ? 1 : 0

  bucket  = google_storage_bucket.bucket.name
  role    = "roles/storage.objectViewer"
  members = var.viewers
}
