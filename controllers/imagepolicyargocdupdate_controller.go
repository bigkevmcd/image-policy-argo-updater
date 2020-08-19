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
	"context"

	imagev1alpha1 "github.com/fluxcd/image-reflector-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1alpha1 "github.com/bigkevmcd/image-policy-argo-updater/api/v1alpha1"
)

const applicationKey = ".spec.application"
const imagePolicyKey = ".spec.imagePolicy"

// ImagePolicyArgoCDUpdateReconciler reconciles a ImagePolicyArgoCDUpdate object
type ImagePolicyArgoCDUpdateReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.bigkevmcd.com,resources=imagepolicyargocdupdates,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.bigkevmcd.com,resources=imagepolicyargocdupdates/status,verbs=get;update;patch

func (r *ImagePolicyArgoCDUpdateReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := newLogger(r.Log.WithValues("imagepolicyargocdupdate", req.NamespacedName))

	var update appsv1alpha1.ImagePolicyArgoCDUpdate
	if err := r.Get(ctx, req.NamespacedName, &update); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	logger.info("loaded the updater", "update", update)

	return ctrl.Result{}, nil
}

func (r *ImagePolicyArgoCDUpdateReconciler) SetupWithManager(mgr ctrl.Manager) error {
	ctx := context.Background()
	// Index the Image Policy (if any) that each ArgoCD Update refers to
	if err := mgr.GetFieldIndexer().IndexField(ctx, &appsv1alpha1.ImagePolicyArgoCDUpdate{}, imagePolicyKey, func(obj runtime.Object) []string {
		updater := obj.(*appsv1alpha1.ImagePolicyArgoCDUpdate)
		return []string{updater.Spec.ImagePolicyRef.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.ImagePolicyArgoCDUpdate{}).
		Watches(&source.Kind{Type: &imagev1alpha1.ImagePolicy{}},
			&handler.EnqueueRequestsFromMapFunc{
				ToRequests: handler.ToRequestsFunc(r.automationsForImagePolicy),
			}).
		Complete(r)
}

// automationsForImagePolicy fetches all the automations that refer to
// a particular ImagePolicy object.
func (r *ImagePolicyArgoCDUpdateReconciler) automationsForImagePolicy(obj handler.MapObject) []ctrl.Request {
	ctx := context.Background()
	var autoList appsv1alpha1.ImagePolicyArgoCDUpdateList
	if err := r.List(ctx, &autoList, client.InNamespace(obj.Meta.GetNamespace()), client.MatchingFields{imagePolicyKey: obj.Meta.GetName()}); err != nil {
		r.Log.Error(err, "failed to list ImageUpdateAutomations for ImagePolicy", "name", types.NamespacedName{
			Name:      obj.Meta.GetName(),
			Namespace: obj.Meta.GetNamespace(),
		})
		return nil
	}
	reqs := make([]ctrl.Request, len(autoList.Items), len(autoList.Items))
	for i := range autoList.Items {
		reqs[i].NamespacedName.Name = autoList.Items[i].GetName()
		reqs[i].NamespacedName.Namespace = autoList.Items[i].GetNamespace()
	}
	return reqs
}
