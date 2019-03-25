# kube-rule

[![Build Status](https://travis-ci.org/chickenzord/kube-rule.svg?branch=master)](https://travis-ci.org/chickenzord/kube-rule)
[![Go Report Card](https://goreportcard.com/badge/github.com/chickenzord/kube-rule)](https://goreportcard.com/report/github.com/chickenzord/kube-rule)
[![codecov](https://codecov.io/gh/chickenzord/kube-rule/branch/master/graph/badge.svg)](https://codecov.io/gh/chickenzord/kube-rule)
[![Automated Docker Build](https://img.shields.io/docker/automated/chickenzord/kube-rule.svg)](https://hub.docker.com/r/chickenzord/kube-rule/)
[![Docker Pulls](https://img.shields.io/docker/pulls/chickenzord/kube-rule.svg)](https://hub.docker.com/r/chickenzord/kube-rule/)

Kubernetes pods admission webhook based on rules CRD. Rewrite of [kube-annotate](https://github.com/chickenzord/kube-annotate) with more generalized uses.


## TL;DR

You will be able to create these Custom Resource Definitions (CRDs) ...

```yaml
apiVersion: kuberule.chickenzord.com/v1alpha1
kind: PodRule
metadata:
  name: staging-rule
  namespace: awesome-staging
spec:
  selector:
    matchLabel:
      tier: app
  mutations:
    annotations:
      example.com/log-enabled: 'true'
      example.com/log-provider: 'a-cheap-logging-stack'
    nodeSelector:
      kubernetes.io/role: app
---
apiVersion: kuberule.chickenzord.com/v1alpha1
kind: PodRule
metadata:
  name: production-rule
  namespace: awesome-production
spec:
  selector:
    matchLabel:
      tier: app
  mutations:
    annotations:
      example.com/log-enabled: 'true'
      example.com/log-provider: 'awesome-and-expensive-logging-stack'
    nodeSelector:
      kubernetes.io/role: app
      example.com/env: production
    tolerations:
    - key: dedicated-env
      operator: Equals
      value: production
```

Don't get it? Basically it allows you to automatically add some predefined specs to selected Pods in certain namespaces. Supports for other resource objects and specs might be added in the future.

## Motivations

> **Why don't you just add those specs to the controller resources directly?** (e.g. `Deployment.spec.template`)

Separation of responsibilities. In a multi-tenants cluster (i.e. shared with multiple teams/organizations) you might not want to complicate already-settled development/deployment flows. With **kube-rule**, cluster operators can define namespace-scoped rules declaratively using CRDs in YAML. 

An example use case: Developers can define how their apps run (e.g. image, command, probes), while cluster operators decide where the apps run (e.g. nodeSelector, affinity, tolerations, etc).

Using CRD simplifies a lot of things. Anyone can leverage any existing tools that understand  how to interact with Kubernetes resources (e.g kubectl, kustomize, ksonnet, etc).

## How to install

### Quick installation

```sh
make quick-install
```

Above command will create CRDs, a namespace `kuberule` and install the controller into it. You might need cluster admin role.

### Recommended installation

We recommend installing kube-rule using Helm Chart (TODO)

## Development

This tool code was bootstrapped using [kubebuilder](http://kubebuilder.netlify.com/) version `1.0.7`.

## TODO

- Helm Chart (high priority)
- Namespace selector in the controller
- Support more specs: tolerations, podAffinity, nodeAffinity, etc
- Support more resources: deployments, statefulsets, daemonsets, etc
- ClusterPodRule CRD (cluster-wide version of PodRule)
