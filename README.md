# image-policy-argo-updater

This is a flux-v2 component for applying updates to ArgoCD applications
automatically, tracking updates to the available image in a repository.

## Installation

Install some prerequisites:

```shell
$ kustomize build github.com/fluxcd/image-reflector-controller/config/default | kubectl apply -f -
```

Then install this controller.

```shell
$ kubectl apply -k config/default
```

## Testing locally

```shell
$ make test
```
