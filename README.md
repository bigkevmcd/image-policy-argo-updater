# image-policy-argo-updater

This is a flux-v2 component for applying updates to ArgoCD applications
automatically, tracking updates to the available image in a repository.

## Testing

Install some prerequisites:

```shell
$ kubectl apply -k github.com/fluxcd/image-reflector-controller/config/default
$ kubectl apply -k github.com/fluxcd/image-automation-controller/config/default
```
