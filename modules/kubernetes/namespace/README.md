# K8S Namespace Module

<!-- NOTE: We use absolute linking here instead of relative linking, because the terraform registry does not support
           relative linking correctly.
-->

This Terraform Module manages Kubernetes
[`Namespaces`](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/).

## What is a Namespace?

A `Namespace` is a Kubernetes resource that can be used to create a virtual environment in your cluster to separate
resources. It allows you to scope resources in your cluster to provide finer grained permission control and resource
quotas.

For example, suppose that you have two different teams managing separate services independently, such that the team
should not be allowed to update or modify the other teams' services. In such a scenario, you would use namespaces to
separate the resources between each team, and implement RBAC roles that only grant access to the namespace if you reside
in the team that manages it.

To summarize, use namespaces to:

- Implement finer grained access control over deployed resources.
- Implement [resource quotas](https://kubernetes.io/docs/concepts/policy/resource-quotas/) to restrict how much of the
  cluster can be utilized by each team.
