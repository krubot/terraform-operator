variable "bindings" {
  description = "Map of role (key) and list of members (value) to add the IAM policies/bindings"
  type        = map(list(string))
}

variable "mode" {
  description = "Mode for adding the IAM policies/bindings, additive and authoritative"
  default     = "additive"
}

variable "entities" {
  description = "Entities list to add the IAM policies/bindings"
  type        = list(string)
}

variable "entity" {
  description = "Entity to add the IAM policies/bindings"
  default     = ""
  type        = string
}
