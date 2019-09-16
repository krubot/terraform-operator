
output "namespace_id" {
  description = "The name of the created namespace."
  value       = element(concat(kubernetes_namespace.namespace.*.id, [""]), 0)
}
