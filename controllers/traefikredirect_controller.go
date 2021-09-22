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
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"time"

	// controller-runtime defines client.Client that reads from Cache write to Kube API serer.
	// defines Manager that manage clients and cache
	// reconciler that compares cluster state, reconcile is called in response to cluster/external events.
	// event fires `reconcile.Request` object as an argument to Reconcile
	// Reconcile returns ctrl.Result
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

	// check for an ExternalName service and create one if it does not exist
	serviceName := fmt.Sprintf(externalServiceNameFormat, traefikRedirect.Name)
	var service corev1.Service
	logger.Info("getting service name:", "service", serviceName)

	err = r.Get(ctx, types.NamespacedName{Namespace: req.Namespace, Name: serviceName}, &service)

	if err != nil {
		if errors.IsNotFound(err) {
			return r.createService(ctx, req, &traefikRedirect, serviceName, externalHost, externalPort)
		} else {
			logger.Error(err, "failed to get service")
			return ctrl.Result{}, err
		}
	}

	logger.Info("managing service:", "service", service)

	// check for Ingress and create one if it does not exist
	ingressName := fmt.Sprintf(ingressNameFormat, traefikRedirect.Name)
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

	logger.Info("managing ingress:", "ingress", ingress)

	logger.Info("updating status:")
	traefikRedirect.Status.LastCheckedTime = &v1.Time{Time: time.Now()}

	if err := r.Status().Update(ctx, &traefikRedirect); err != nil {
		logger.Error(err, "failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
// Manager manages shared dependencies such as Cache and Client
func (r *TraefikRedirectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&traefikv1.TraefikRedirect{}).
		// will reconcile the service if it's modified/deleted externally
		Owns(&corev1.Service{}).
		// will reconcile the ingress if it's modified/deleted externally
		Owns(&networkv1.Ingress{}).
		Complete(r)
}
