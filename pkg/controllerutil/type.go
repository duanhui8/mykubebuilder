package controllerutil

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
)

type RequestCtx struct {
	Ctx      context.Context
	Req      ctrl.Request
	Log      logr.Logger
	Recorder record.EventRecorder
}
