/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"path/filepath"
	"testing"

	"context"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	argov1alpha1 "github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	appsv1alpha1 "github.com/bigkevmcd/image-policy-argo-updater/api/v1alpha1"
	imagev1alpha1 "github.com/fluxcd/image-reflector-controller/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	cfg        *rest.Config
	k8sClient  client.Client
	testEnv    *envtest.Environment
	k8sManager ctrl.Manager
)

const (
	updaterName      = "test-updater"
	updaterNamespace = "updaters"
	timeout          = 2 * time.Second
	argoAppName      = "my-demo-app"
	argoAppNamespace = "argocd"
	policyName       = "my-policy"
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "config", "crd", "bases"),
			filepath.Join("testdata", "crds"),
		},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	schemes := []func(*runtime.Scheme) error{
		appsv1alpha1.AddToScheme, argov1alpha1.AddToScheme, imagev1alpha1.AddToScheme,
	}
	for _, v := range schemes {
		err = v(scheme.Scheme)
		Expect(err).NotTo(HaveOccurred())
	}

	// +kubebuilder:scaffold:scheme

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	err = (&ImagePolicyArgoCDUpdateReconciler{
		Client: k8sManager.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("ImagePolicyArgoCDUpdateReconciler"),
		Scheme: scheme.Scheme,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctrl.SetupSignalHandler())
		Expect(err).ToNot(HaveOccurred())
	}()

	k8sClient = k8sManager.GetClient()
	Expect(k8sClient).ToNot(BeNil())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

var _ = Describe("ImagePolicyArgoCDUpdate", func() {
	var (
		updater     *appsv1alpha1.ImagePolicyArgoCDUpdate
		argoApp     *argov1alpha1.Application
		policy      *imagev1alpha1.ImagePolicy
		latestImage string
	)

	BeforeEach(func() {
		updater = &appsv1alpha1.ImagePolicyArgoCDUpdate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      updaterName,
				Namespace: updaterNamespace,
			},
			Spec: appsv1alpha1.ImagePolicyArgoCDUpdateSpec{
				ApplicationRef: corev1.ObjectReference{
					Name:      argoAppName,
					Namespace: argoAppNamespace,
				},
				ImagePolicyRef: corev1.LocalObjectReference{
					Name: policyName,
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), updater)).To(Succeed())

		argoApp = &argov1alpha1.Application{
			ObjectMeta: metav1.ObjectMeta{
				Name:      argoAppName,
				Namespace: argoAppNamespace,
			},
			Spec: argov1alpha1.ApplicationSpec{
				Source: argov1alpha1.ApplicationSource{},
			},
		}
		Expect(k8sClient.Create(context.Background(), argoApp)).To(Succeed())

		policy = &imagev1alpha1.ImagePolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      policyName,
				Namespace: updaterNamespace,
			},
			Spec: imagev1alpha1.ImagePolicySpec{
				ImageRepositoryRef: corev1.LocalObjectReference{
					Name: "not-used",
				},
				Policy: imagev1alpha1.ImagePolicyChoice{
					SemVer: &imagev1alpha1.SemVerPolicy{
						Range: "1.14.x",
					},
				},
			},
		}
		Expect(k8sClient.Create(context.Background(), policy)).To(Succeed())
	})

	JustBeforeEach(func() {
		ctx := context.Background()
		policy.Status = imagev1alpha1.ImagePolicyStatus{
			LatestImage: latestImage,
		}
		Expect(k8sClient.Status().Update(ctx, policy)).To(Succeed())
	})

	AfterEach(func() {
		ctx := context.Background()
		Expect(k8sClient.Delete(ctx, updater)).To(Succeed())
		Expect(k8sClient.Delete(ctx, argoApp)).To(Succeed())
		Expect(k8sClient.Delete(ctx, policy)).To(Succeed())
	})

	Context("an ImagePolicy is updated", func() {
		Context("associated with a ImagePolicyArgoCDUpdate", func() {
			BeforeEach(func() {
				latestImage = "1.14.5"
			})
			It("triggers the update of the related ArgoCD application", func() {
				Eventually(func() argov1alpha1.KustomizeImages {
					loaded := &argov1alpha1.Application{}
					Expect(k8sClient.Get(context.Background(), types.NamespacedName{
						Name:      argoAppName,
						Namespace: argoAppNamespace,
					}, loaded)).NotTo(HaveOccurred())
					if loaded.Spec.Source.Kustomize != nil {
						return loaded.Spec.Source.Kustomize.Images
					}
					return argov1alpha1.KustomizeImages{}
				}, timeout, time.Millisecond*500).Should(Equal(argov1alpha1.KustomizeImages{argov1alpha1.KustomizeImage(latestImage)}))
			})
		})

		Context("not associated with a ImagePolicyArgoCDUpdate", func() {
		})
	})
})
