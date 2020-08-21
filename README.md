# image-policy-argo-updater

This is a flux-v2 component for applying updates to ArgoCD applications
automatically, tracking updates to the available image in a repository.

## Installation

Install some prerequisites:

```shell
$ kustomize build github.com/fluxcd/image-reflector-controller/config/default | kubectl apply -f -
```

## Testing locally

```shell
$ make test
```
