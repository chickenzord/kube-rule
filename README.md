# kube-rule

[![Build Status](https://travis-ci.org/chickenzord/kube-rule.svg?branch=master)](https://travis-ci.org/chickenzord/kube-rule)
[![Go Report Card](https://goreportcard.com/badge/github.com/chickenzord/kube-rule)](https://goreportcard.com/report/github.com/chickenzord/kube-rule)
[![codecov](https://codecov.io/gh/chickenzord/kube-rule/branch/master/graph/badge.svg)](https://codecov.io/gh/chickenzord/kube-rule)
[![Automated Docker Build](https://img.shields.io/docker/automated/chickenzord/kube-rule.svg)](https://hub.docker.com/r/chickenzord/kube-rule/)
[![Docker Pulls](https://img.shields.io/docker/pulls/chickenzord/kube-rule.svg)](https://hub.docker.com/r/chickenzord/kube-rule/)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fchickenzord%2Fkube-rule.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fchickenzord%2Fkube-rule?ref=badge_shield)

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
    imagePullSecrets:
    - name: dockerhub-creds
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
    imagePullSecrets:
    - name: dockerhub-creds
```

Don't get it? Basically it allows you to automatically add some predefined specs to selected Pods in certain namespaces. Supports for other resource objects and specs might be added in the future.

## Motivations

> **Why don't you just add those specs to the controller resources directly?** (e.g. `Deployment.spec.template`)

Separation of responsibilities. In a multi-tenants cluster (i.e. shared with multiple teams/organizations) you might not want to complicate already-settled development/deployment flows. With **kube-rule**, Developers can define how their apps run (e.g. image, command, probes), while cluster operators decide where the apps run (e.g. nodeSelector, affinity, tolerations, etc).

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

## Goals and Non-Goals

**What mutation spec should kube-rule supports?**

Basically kube-rule aims to support Pod specs that can be implicitly **decided by cluster admins** depending on where the app runs (i.e. environment-dependent).

Example of specs mutation that **will not get supported** by kube-rule:

- `containers.image`: Images deployed should be decided by CI/CD in app layer tooling.
- `containers.commands`/`containers.args`: Overriding them requires knowledge of the container image used.
- `volumes`/`volumeMounts`: Volumes used are closely tied to app logic.
- etc...

Feel free to request more specs by describing your valid use-case in Pull Requests.

## TODO

- Helm Chart (high priority)
- Namespace selector in the controller
- Support more specs: 
  - ~~tolerations~~
  - ~~podAffinity~~
  - ~~nodeAffinity~~
  - containers.resources
  - etc
- Support more resources: deployments, statefulsets, daemonsets, etc
- ClusterPodRule CRD (cluster-wide version of PodRule)


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fchickenzord%2Fkube-rule.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fchickenzord%2Fkube-rule?ref=badge_large)