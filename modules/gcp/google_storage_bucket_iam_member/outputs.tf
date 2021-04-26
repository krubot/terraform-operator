# ---------------------------------------------------------------------------------------------------------------------
# GCS BUCKET IAM OUTPUTS
# These are the variable outputs of the gcs bucket IAM creation
# ---------------------------------------------------------------------------------------------------------------------

output "bindings_by_member" {
  value       = local.bindings_by_member
  description = "List of bindings for entities unwinded by members."
}
