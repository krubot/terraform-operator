# ---------------------------------------------------------------------------------------------------------------------
# GCS BUCKET OUTPUTS
# These are the variable outputs of the gcs bucket creation
# ---------------------------------------------------------------------------------------------------------------------

output "name" {
  description = "Bucket name (for single use)."
  value       = google_storage_bucket.bucket.name
}

output "url" {
  description = "Bucket URL (for single use)."
  value       = google_storage_bucket.bucket.url
}
