# Terraform Operator

Currently still in POC stage, there are plans to extend this functionality soon. Feel free to create issues but note it's only one person working here 😄

Simply want the image run `make image`.

## Pre-requisites

The following are needed to run this repo:

 - Kubernetes cluster - [minikube](https://github.com/kubernetes/minikube) is a simple tool for this
 - Operator SDK - install guide can be found [here](https://github.com/operator-framework/operator-sdk)
 - Helm - binary can be found [here](https://github.com/helm/helm)
 - Quay.io login - Login page can be found [here](https://quay.io/)


Make sure to have also forked and cloned the repo if you are deploying with [flux](https://github.com/weaveworks/flux). This is recommended since it'll sync all your changes and help you to deploy consistently.

## Building the operator

To build the operator image run the following:

```sh
operator-sdk build quay.io/YOURUSER/terraform-operator:latest
```

Now push to `quay.io` by running:

```sh
docker push quay.io/YOURUSER/terraform-operator:latest
```

## Running Helm and Flux

To run the pipeline all the way through please deploy `helm` and `flux` with the following command:

```sh
kubectl apply -f manifest/helm.yaml
kubectl apply -f manifest/flux.yaml
```

Calico is also here if you deploying kubernetes from scratch. Apply with a similar command.

## Running some tests

To test that the deployment please checkout the `infra` namespace and validate in the logs that the terraform operator is running correctly.

## Workflow identity service account

If running this on the google kubernetes engine then make sure you have workload identity enable. The link below is to the terraform config argument where this must be set:

https://www.terraform.io/docs/providers/google/r/container_cluster.html#workload_identity_config

The using the gcloud cli you can generate the `terraform-operator` service account and permissions:

```sh
$ gcloud --project=<project> iam service-accounts create terraform-operator --display-name "Terraform operator service account"
$ gcloud --project=<project> iam service-accounts add-iam-policy-binding --role "roles/iam.workloadIdentityUser" --member "serviceAccount:<project>.svc.id.goog[infra/terraform-operator]" terraform-operator@<project>.iam.gserviceaccount.com
$ gcloud projects add-iam-policy-binding <project> --member='serviceAccount:terraform-operator@<project>.iam.gserviceaccount.com' --role='roles/cloudsql.admin'
```
(`<project>` is the gcp project id)

This can be then added to the release values and used in the helm deploy:

```yaml
serviceAccount:
  create: true
  name: terraform-operator
  gcpServiceAccount:
    create: true
    name: terraform-operator@<project>.iam.gserviceaccount.com
```
