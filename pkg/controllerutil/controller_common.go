package controllerutil

import (
	"strings"

	"duangh/mykubebuilder/pkg/constant"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// CheckedRequeueWithError passes the error through to the controller
// manager, it ignores unknown errors.
func CheckedRequeueWithError(err error, logger logr.Logger, msg string, keysAndValues ...interface{}) (reconcile.Result, error) {
	if apierrors.IsNotFound(err) {
		return Reconciled()
	}
	return RequeueWithError(err, logger, msg, keysAndValues...)
}

// Reconciled returns an empty result with nil error to signal a successful reconcile
// to the controller manager
func Reconciled() (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

// RequeueWithError requeues when an error occurs
func RequeueWithError(err error, logger logr.Logger, msg string, keysAndValues ...interface{}) (reconcile.Result, error) {
	if msg == "" {
		logger.Info(err.Error())
	} else {
		// Info log the error message and then let the reconciler dump the stacktrace
		logger.Info(msg, keysAndValues...)
	}
	return reconcile.Result{}, err
}

// RecordCreatedEvent records an event when a CR created successfully
func RecordCreatedEvent(r record.EventRecorder, cr client.Object) {
	if r != nil && cr.GetGeneration() == 1 {
		r.Eventf(cr, corev1.EventTypeNormal, constant.ReasonCreatedCR, "Created %s: %s", strings.ToLower(cr.GetObjectKind().GroupVersionKind().Kind), cr.GetName())
	}
}
