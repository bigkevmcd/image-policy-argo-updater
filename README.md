# image-policy-argo-updater

This is a flux-v2 component for applying updates to ArgoCD applications
automatically, tracking updates to the available image in a repository.

## Installation

Install some prerequisites:

```shell
$ kubectl apply -k github.com/fluxcd/image-reflector-controller/config/default
```

## Testing locally

```shell
$ make test
```
