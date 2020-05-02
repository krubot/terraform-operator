
# ---------------------------------------------------------------------------------------------------------------------
# GCS BUCKET IAM BINDING CONFIGURATION
# ---------------------------------------------------------------------------------------------------------------------

resource "google_storage_bucket_iam_binding" "storage_bucket_iam_authoritative" {
  for_each = local.set_authoritative
  bucket   = local.bindings_authoritative[each.key].name
  role     = local.bindings_authoritative[each.key].role
  members  = local.bindings_authoritative[each.key].members
}

resource "google_storage_bucket_iam_member" "storage_bucket_iam_additive" {
  for_each = local.set_additive
  bucket   = local.bindings_additive[each.key].name
  role     = local.bindings_additive[each.key].role
  member   = local.bindings_additive[each.key].member
}
