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

 Section can be skipped if you have already setup `helm tiller` and don't wish to run flux. Firstly create a service account and a cluster role binding for Tiller:

 ```sh
 kubectl -n kube-system create sa tiller

 kubectl create clusterrolebinding tiller-cluster-rule \
     --clusterrole=cluster-admin \
     --serviceaccount=kube-system:tiller
 ```

 Then deploy Tiller in `kube-system` namespace:

 ```sh
 helm init --skip-refresh --upgrade --wait --service-account tiller
 ```

 Add the Flux repository of Weaveworks:

 ```sh
 helm repo add fluxcd https://charts.fluxcd.io
 ```

Apply the Helm Release CRD:

 ```sh
 kubectl apply -f https://raw.githubusercontent.com/fluxcd/flux/helm-0.10.1/deploy-helm/flux-helm-release-crd.yaml
 ```

Next install flux and replace the `git.url` with your repos url:


```sh
helm upgrade -i flux \
--set helmOperator.create=true \
--set helmOperator.createCRD=false \
--set git.url=git@github.com:YOURUSER/terraform-operator \
--namespace flux \
fluxcd/flux
```

Next obtain flux's public ssh-key by running:

```sh
kubectl -n flux logs deployment/flux | grep identity.pub | cut -d '"' -f2
```

In order to sync your cluster state with git you need to copy the public key and create a deploy key with write access on your GitHub repository.

Open GitHub, navigate to your fork, go to **Setting > Deploy keys**, click on **Add deploy key**, give it a name, check **Allow write access**, paste the Flux public key and click **Add key**.

Once Flux has confirmed access to the repository, it will start deploying the workload. After a while you will be able to see the Helm releases listed like so:

## Running some tests

To test that the deployment has worked you can run some custom resource files with the follow command:

```sh
kubectl -n infra apply -f https://raw.githubusercontent.com/krubot/terraform-operator/master/deploy/crds/terraform_v1alpha1_provider_cr.yaml
```
