# Terraform Operator

Currently still in POC stage, there are plans to extend this functionality soon. Feel free to create issues but note it's only one person working here 😄

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
 helm init --skip-refresh --upgrade --service-account tiller
 ```

 Add the Flux repository of Weaveworks:

 ```sh
 helm repo add weaveworks https://weaveworks.github.io/flux
 ```

Next install flux and replace the `git.url` with your repos url:


```sh
helm install --name flux \
--set helmOperator.create=true \
--set git.url=ssh://git@github.com/YOURUSER/terraform-operator \
--set helmOperator.git.chartsPath=charts \
--namespace flux \
weaveworks/flux
```

Next obtain flux's public ssh-key by running:

```sh
FLUX_POD=$(kubectl get pods --namespace flux -l "app=flux,release=flux" -o jsonpath="{.items[0].metadata.name}")
kubectl -n flux logs $FLUX_POD | grep identity.pub | cut -d '"' -f2
```

In order to sync your cluster state with git you need to copy the public key and create a deploy key with write access on your GitHub repository.

Open GitHub, navigate to your fork, go to **Setting > Deploy keys**, click on **Add deploy key**, give it a name, check **Allow write access**, paste the Flux public key and click **Add key**.

Once Flux has confirmed access to the repository, it will start deploying the workload. After a while you will be able to see the Helm releases listed like so:

## Running some tests

To test that the deployment has worked you can run some custom resource files in the testing folder.
