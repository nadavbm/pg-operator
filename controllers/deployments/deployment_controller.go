/*
Copyright 2023.

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

package deployments

import (
	"context"
	"time"

	"github.com/nadavbm/zlog"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	deploymentsv1alpha1 "example.com/pg/apis/deployments/v1alpha1"
	"example.com/pg/controllers/specs"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	Logger *zlog.Logger
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=deployments.example.com,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=deployments.example.com,resources=deployments/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=deployments.example.com,resources=deployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Deployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := zlog.New()
	r.Logger = logger

	r.Logger.Info("Start reconcile", zap.String("namespace", req.NamespacedName.Namespace))

	var resource deploymentsv1alpha1.Deployment
	if err := r.Client.Get(context.Background(), req.NamespacedName, &resource); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("resource is not found, probably deleted. skipping..", zap.String("namespace", req.Namespace))
			return ctrl.Result{Requeue: false, RequeueAfter: 0}, nil
		}
		r.Logger.Error("could not fetch resource", zap.String("type", resource.Kind))
		return ctrl.Result{Requeue: true, RequeueAfter: time.Minute}, err
	}

	var object appsv1.Deployment
	if err := r.Get(ctx, req.NamespacedName, &object); err != nil {
		if errors.IsNotFound(err) {
			r.Logger.Info("create deployment", zap.String("namespace", req.Namespace))
			obj := specs.BuildDeployment(req.Namespace, &resource)
			if err := r.Create(ctx, obj); err != nil {
				r.Logger.Error("could not create", zap.String("object kind", obj.Kind))
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		if err := r.Update(ctx, &object); err != nil {
			if errors.IsInvalid(err) {
				r.Logger.Error("invalid update", zap.String("object", object.Name))
			} else {
				r.Logger.Error("unable to update", zap.String("object", object.Name))
			}
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&deploymentsv1alpha1.Deployment{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
