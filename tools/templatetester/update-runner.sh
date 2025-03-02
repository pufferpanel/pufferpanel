#! /bin/sh

helm upgrade --namespace arc-runners -f docker.yaml  -f template-tester.yaml template-tester oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set