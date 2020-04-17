
# ---------------------------------------------------------------------------------------------------------------------
# REQUIRED PARAMETERS
# These variables are expected to be passed in by the operator
# ---------------------------------------------------------------------------------------------------------------------

variable "project" {
  description = "Bucket project id."
  type        = string
}

variable "name" {
  description = "Bucket name suffixes."
  type        = string
}

# ---------------------------------------------------------------------------------------------------------------------
# OPTIONAL PARAMETERS
# Generally, these values won't need to be changed.
# ---------------------------------------------------------------------------------------------------------------------

variable "location" {
  description = "Bucket location."
  type        = string
  default     = "EU"
}

variable "storage_class" {
  description = "Bucket storage class."
  type        = string
  default     = "MULTI_REGIONAL"
}

variable "force_destroy" {
  description = "Optional map of lowercase unprefixed name => boolean, defaults to false."
  default     = {}
}

variable "versioning" {
  description = "Optional map of lowercase unprefixed name => boolean, defaults to false."
  default     = {}
}

variable "encryption_key_names" {
  description = "Optional map of lowercase unprefixed name => string, empty strings are ignored."
  default     = {}
}

variable "bucket_policy_only" {
  description = "Disable ad-hoc ACLs on specified buckets. Defaults to true. Map of lowercase unprefixed name => boolean"
  default     = {}
}

variable "lifecycle_rules" {
  type = set(object({
    action    = map(string)
    condition = map(string)
  }))
  default     = []
  description = "List of lifecycle rules to configure. Format is the same as described in provider documentation https://www.terraform.io/docs/providers/google/r/storage_bucket.html#lifecycle_rule except condition.matches_storage_class should be a comma delimited string."
}

variable "admins" {
  description = "IAM-style members who will be granted roles/storage.objectAdmin on all buckets."
  type        = list(string)
  default     = []
}

variable "creators" {
  description = "IAM-style members who will be granted roles/storage.objectCreators on all buckets."
  type        = list(string)
  default     = []
}

variable "viewers" {
  description = "IAM-style members who will be granted roles/storage.objectViewer on all buckets."
  type        = list(string)
  default     = []
}

variable "labels" {
  description = "Labels to be attached to the buckets"
  default     = {}
}
