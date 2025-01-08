/*
Copyright 2025.

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

package controller

import (
	"context"
	"slices"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	duanhui8v1 "duangh/mykubebuilder/api/v1"

	intctrlutil "duangh/mykubebuilder/pkg/controllerutil"
)

// ClusterDefinitionReconciler reconciles a ClusterDefinition object
type ClusterDefinitionReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var clusterDefUpdateHandlers = map[string]func(client client.Client, ctx context.Context, clusterDef *duanhui8v1.ClusterDefinition) error{}

//+kubebuilder:rbac:groups=duanhui8.duangh,resources=clusterdefinitions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=duanhui8.duangh,resources=clusterdefinitions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=duanhui8.duangh,resources=clusterdefinitions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterDefinition object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ClusterDefinitionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqCtx := intctrlutil.RequestCtx{
		Ctx:      ctx,
		Req:      req,
		Log:      log.FromContext(ctx).WithValues("clusterDefinition", req.NamespacedName),
		Recorder: r.Recorder,
	}
	clusterDef := &duanhui8v1.ClusterDefinition{}
	if err := r.Client.Get(reqCtx.Ctx, reqCtx.Req.NamespacedName, clusterDef); err != nil {
		return intctrlutil.CheckedRequeueWithError(err, reqCtx.Log, "")
	}

	if clusterDef.Status.ObservedGeneration == clusterDef.Generation && slices.Contains(clusterDef.Status.GetTerminalPhases(), clusterDef.Status.Phase) {
		return intctrlutil.Reconciled()
	}
	r.reconcile(reqCtx, clusterDef)
	intctrlutil.RecordCreatedEvent(r.Recorder, clusterDef)
	return ctrl.Result{}, nil
}

func (r *ClusterDefinitionReconciler) reconcile(rctx intctrlutil.RequestCtx, clusterDef *duanhui8v1.ClusterDefinition) (*ctrl.Result, error) {
	if err := r.reconcileTopologies(rctx, clusterDef); err != nil {
		res, err1 := intctrlutil.CheckedRequeueWithError(err, rctx.Log, "")
		return &res, err1
	}
	for _, handler := range clusterDefUpdateHandlers {
		if err := handler(r.Client, rctx.Ctx, clusterDef); err != nil {
			res, err1 := intctrlutil.CheckedRequeueWithError(err, rctx.Log, "")
			return &res, err1
		}
	}
	return nil, nil
}
func (r *ClusterDefinitionReconciler) reconcileTopologies(rctx intctrlutil.RequestCtx, clusterDef *duanhui8v1.ClusterDefinition) error {
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterDefinitionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&duanhui8v1.ClusterDefinition{}).
		Complete(r)
}
