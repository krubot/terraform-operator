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
