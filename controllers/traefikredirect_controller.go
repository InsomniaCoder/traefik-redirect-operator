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
	"fmt"

	traefikv1 "github.com/InsomniaCoder/traefik-redirect-operator/api/v1"
	_ "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	ingressNameFormat         string = "%s-traefik-ingress"
	externalServiceNameFormat string = "%s-svc-external"
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
	logger := log.FromContext(ctx)
	logger.Info("Start reconciling....")

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

	logger.Info("got Traefik:", "traefik", traefikRedirect)

	traefikType := traefikRedirect.Spec.TraefikType
	traefikHost := traefikRedirect.Spec.TraefikHost
	externalHost := traefikRedirect.Spec.RedirectTo
	externalPort := traefikRedirect.Spec.Port

	logger.Info("traefik type: ", "type", traefikType)
	logger.Info("traefik host: ", "host", traefikHost)
	logger.Info("external host: ", "host", externalHost)
	logger.Info("external port: ", "port", externalPort)

	logger.Info("target namespace: ", "ns", req.Namespace)

	// check for Ingress and create one if it does not exist
	ingressName := fmt.Sprintf(ingressNameFormat, traefikRedirect.Name)
	serviceName := fmt.Sprintf(externalServiceNameFormat, traefikRedirect.Name)
	var ingress networkv1.Ingress
	err = r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: ingressName}, &ingress)

	if err != nil {
		if errors.IsNotFound(err) {
			return r.createIngress(ctx, req, &traefikRedirect, ingressName, serviceName, traefikType, traefikHost)
		} else {
			logger.Error(err, "failed to get Ingress")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TraefikRedirectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&traefikv1.TraefikRedirect{}).
		Complete(r)
}
