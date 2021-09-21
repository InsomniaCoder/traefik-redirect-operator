/*
Copyright 2021.

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

	traefikv1 "github.com/InsomniaCoder/traefik-redirect-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
//externalServiceNameFormat string = "%s-external"
)

// TraefikRedirectReconciler reconciles a TraefikRedirect object
type TraefikRedirectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=traefik.porpaul,resources=traefikredirects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=service,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extersions/v1beta1,resources=ingress,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=traefik.porpaul,resources=traefikredirects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=traefik.porpaul,resources=traefikredirects/finalizers,verbs=update

func (r *TraefikRedirectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.Log.WithValues("TraefikRedirect", req.NamespacedName)
	logger.Info("start reconciling")

	var traefikRedirect traefikv1.TraefikRedirect
	err := r.Get(ctx, req.NamespacedName, &traefikRedirect)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			logger.Info("traefik redirect resource not found.")
			return ctrl.Result{}, nil
		}

		logger.Error(err, "failed to get traefik redirect")
		return ctrl.Result{}, err
	}

	logger.Info("got Traefik object: %v", traefikRedirect)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TraefikRedirectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&traefikv1.TraefikRedirect{}).
		Complete(r)
}
