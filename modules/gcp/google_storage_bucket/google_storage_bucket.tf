
# ---------------------------------------------------------------------------------------------------------------------
# GCS BUCKET CONFIGURATION
# ---------------------------------------------------------------------------------------------------------------------

resource "google_storage_bucket" "bucket" {
  name          = var.name
  location      = var.location
  storage_class = var.storage_class
  labels        = merge(var.labels, { name = replace("${lower(var.name)}", ".", "-") })

  force_destroy      = lookup(var.force_destroy,lower(var.name),false)
  bucket_policy_only = lookup(var.bucket_policy_only,lower(var.name),true)

  versioning {
    enabled = lookup(var.versioning,lower(var.name),false)
  }

  # Having a permanent encryption block with default_kms_key_name = "" works but results in terraform applying a change every run
  # There is no enabled = false attribute available to ask terraform to ignore the block
  dynamic "encryption" {
    # If an encryption key name is set for this bucket name -> Create a single encryption block
    for_each = trimspace(lookup(var.encryption_key_names, lower(var.name), "")) != "" ? [true] : []
    content {
      default_kms_key_name = trimspace(
        lookup(
          var.encryption_key_names,
          lower(var.name),
          "Error retrieving kms key name", # Should be unreachable due to the for_each check
          # Omitting default is deprecated & can help show if there was a bug
          # https://www.terraform.io/docs/configuration/functions/lookup.html
        )
      )
    }
  }

  dynamic "lifecycle_rule" {
    for_each = var.lifecycle_rules
    content {
      action {
        type          = lifecycle_rule.value.action.type
        storage_class = lookup(lifecycle_rule.value.action, "storage_class", null)
      }
      condition {
        age                   = lookup(lifecycle_rule.value.condition, "age", null)
        created_before        = lookup(lifecycle_rule.value.condition, "created_before", null)
        with_state            = lookup(lifecycle_rule.value.condition, "with_state", null)
        matches_storage_class = contains(keys(lifecycle_rule.value.condition), "matches_storage_class") ? split(",", lifecycle_rule.value.condition["matches_storage_class"]) : null
        num_newer_versions    = lookup(lifecycle_rule.value.condition, "num_newer_versions", null)
      }
    }
  }
}
