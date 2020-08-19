module github.com/bigkevmcd/image-policy-argo-updater

go 1.13

require (
	github.com/fluxcd/image-reflector-controller v0.0.0-20200819120130-b302367aac9e
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.12.1
	github.com/onsi/gomega v1.10.1
	k8s.io/api v0.18.6
	k8s.io/apimachinery v0.18.6
	k8s.io/client-go v0.18.6
	sigs.k8s.io/controller-runtime v0.6.2
)
