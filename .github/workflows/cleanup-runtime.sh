#!/usr/bin/env bash

set -o nounset

readonly version=$(cat VERSION)
readonly git_sha=$(git rev-parse HEAD)
readonly git_timestamp=$(TZ=UTC git show --quiet --date='format-local:%Y%m%d%H%M%S' --format="%cd")
readonly slug=${version}-${git_timestamp}-${git_sha:0:16}

readonly tiller_service_account=tiller
readonly tiller_namespace=kube-system

echo "Cleanup riff Build"
kubectl delete riff -n $NAMESPACE --all
kubectl delete -f https://storage.googleapis.com/projectriff/riff-system/snapshots/riff-build-${slug}.yaml
kubectl delete -f https://storage.googleapis.com/projectriff/riff-buildtemplate/riff-application-clusterbuildtemplate.yaml
kubectl delete -f https://storage.googleapis.com/projectriff/riff-buildtemplate/riff-function-clusterbuildtemplate.yaml

echo "Cleanup Knative Build"
kubectl delete -f https://storage.googleapis.com/knative-releases/build/previous/v0.7.0/build.yaml

if [ $RUNTIME = "core" ]; then
  echo "Cleanup riff Core Runtime"
  kubectl delete -f https://storage.googleapis.com/projectriff/riff-system/snapshots/riff-core-${slug}.yaml

elif [ $RUNTIME = "knative" ]; then
  echo "Cleanup Istio"
  helm delete --purge istio
  kubectl delete namespace istio-system
  kubectl get customresourcedefinitions.apiextensions.k8s.io -oname | grep istio.io | xargs -L1 kubectl delete

  echo "Cleanup Knative Serving"
  kubectl delete knative -n $NAMESPACE --all
  retry kubectl apply -f https://storage.googleapis.com/knative-releases/serving/previous/v0.9.0/serving-post-1.14.yaml

  echo "Cleanup riff Knative Runtime"
  kubectl delete -f https://storage.googleapis.com/projectriff/riff-system/snapshots/riff-knative-${slug}.yaml

fi

echo "Remove Helm"
helm reset

kubectl delete serviceaccount ${tiller_service_account} -n ${tiller_namespace}
kubectl delete clusterrolebinding "${tiller_service_account}-cluster-admin"
