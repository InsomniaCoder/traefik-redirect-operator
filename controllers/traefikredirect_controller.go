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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TraefikRedirectReconciler reconciles a TraefikRedirect object
type TraefikRedirectReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=traefik.porpaul,resources=traefikredirects,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=v1,resources=service,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=extersions/v1beta1,resources=ingress,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=traefik.porpaul,resources=traefikredirects/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=traefik.porpaul,resources=traefikredirects/finalizers,verbs=update

func (r *TraefikRedirectReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := ctrl.Log.WithValues("TraefikRedirect", req.NamespacedName)
	l.Info("start reconciling")


	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TraefikRedirectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&traefikv1.TraefikRedirect{}).
		Complete(r)
}
