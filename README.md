# Terraform Operator

Currently still in POC stage, there are plans to extend this functionality soon. Feel free to create issues but note it's only one person working here 😄

Simply want the binary run `make`.

## Pre-requisites

The following are needed to run this repo:

 - Kubernetes cluster - [minikube](https://github.com/kubernetes/minikube) is a simple tool for this
 - Helm - binary can be found [here](https://github.com/helm/helm)

Make sure to have also forked and cloned the repo if you are deploying with [flux](https://github.com/weaveworks/flux). This is recommended since it'll sync all your changes and help you to deploy consistently.

## Building the operator

To build the operator image run the following:

```sh
IMG=<image-repo> make docker-build
```

Now push to your repo run:

```sh
IMG=<image-repo> make docker-push
```

## Running Helm and Flux

To install helm flux please run the following:

```sh
kubectl create ns flux
```

Next we need to create the CRD's for helm operator in advance of creating the deployment:

```sh
kubectl apply -f https://raw.githubusercontent.com/fluxcd/helm-operator/1.1.0/deploy/crds.yaml
```

Now we can add the fluxcd charts and run a install:

```sh
helm repo add fluxcd https://charts.fluxcd.io

helm upgrade -i flux fluxcd/flux \
    --namespace flux \
    --set git.url=git@github.com:krubot/terraform-operator \
    --set git.readonly=true \
    --set git.path=deploy \
    --set rbac.pspEnabled=true

helm upgrade -i helm-operator fluxcd/helm-operator \
    --namespace flux \
    --set git.ssh.secretName=flux-git-deploy \
    --set helm.versions=v3
```

The following need to now be run to get the pubic ssh key:

```sh
kubectl -n flux logs deployment/flux | grep identity.pub | cut -d '"' -f2
```

This should output should be the whole key which you add to your deployments configuration within your github repo. This key does not need write access so don't tick this box.

## Running some tests

To test that the deployment please checkout the `infra` namespace and validate in the logs that the terraform operator is running correctly.

## Workflow identity service account

If running this on the google kubernetes engine then make sure you have workload identity enable. The link below is to the terraform config argument where this must be set:

https://www.terraform.io/docs/providers/google/r/container_cluster.html#workload_identity_config

The using the gcloud cli you can generate the `terraform-operator` service account and permissions:

```sh
$ gcloud --project=<project> iam service-accounts create terraform-operator --display-name "Terraform operator service account"
$ gcloud --project=<project> iam service-accounts add-iam-policy-binding --role "roles/iam.workloadIdentityUser" --member "serviceAccount:<project>.svc.id.goog[infra/terraform-operator]" terraform-operator@<project>.iam.gserviceaccount.com
$ gcloud projects add-iam-policy-binding <project> --member='serviceAccount:terraform-operator@<project>.iam.gserviceaccount.com' --role='roles/storage.admin'
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
