package controllers

import (
	"context"
	traefikv1 "github.com/InsomniaCoder/traefik-redirect-operator/api/v1"
	networkv1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *TraefikRedirectReconciler) createIngress(ctx context.Context, req ctrl.Request, tr *traefikv1.TraefikRedirect, ingressName string, serviceName string, traefikType traefikv1.TraefikType, host string) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	ingress := &networkv1.Ingress{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      ingressName,
			Namespace: req.Namespace,
			Annotations: map[string]string{
				"kubernetes.io/ingress.class": "traefik",
			},
			Labels: map[string]string{
				"traffic-type": string(traefikType),
			},
		},
		Spec: networkv1.IngressSpec{
			Rules: []networkv1.IngressRule{{
				Host: host,
				IngressRuleValue: networkv1.IngressRuleValue{
					HTTP: &networkv1.HTTPIngressRuleValue{
						Paths: []networkv1.HTTPIngressPath{
							{
								Path: "/",
								Backend: networkv1.IngressBackend{
									Service: &networkv1.IngressServiceBackend{
										Name: serviceName,
										Port: networkv1.ServiceBackendPort{
											Name: "http",
										},
									},
								},
							},
						},
					},
				},
			},
			},
		},
	}
	// Set Traefik redirect crd instance as the owner and controller
	err := ctrl.SetControllerReference(tr, ingress, r.Scheme)
	if err != nil {
		logger.Error(err, "error setting controller reference", "namespace", req.Namespace, "name", ingressName)
		return ctrl.Result{}, err
	}

	logger.Info("creating ingress", "ingress", ingressName)

	err = r.Create(ctx, ingress)

	if err != nil {
		logger.Error(err, "failed to create ingress", "ingress", ingressName)
		return ctrl.Result{}, err
	}

	logger.Info("ingress created successfully", "ingress", ingress)

	// return response for created successfully
	return ctrl.Result{Requeue: true}, nil
}
