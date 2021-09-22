package controllers

import (
	"context"
	traefikv1 "github.com/InsomniaCoder/traefik-redirect-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func (r *TraefikRedirectReconciler) createService(ctx context.Context, req ctrl.Request, tr *traefikv1.TraefikRedirect, serviceName string, externalHost string, port int) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	svc := &corev1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      serviceName,
			Namespace: req.Namespace,
		},
		Spec: corev1.ServiceSpec{
			ExternalName: externalHost,
			Type:         corev1.ServiceType("ExternalName"),
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.Protocol("TCP"),
					TargetPort: intstr.FromInt(port),
					Port:       int32(port),
				},
			},
		},
	}

	// Set Traefik redirect crd instance as the owner and controller
	err := ctrl.SetControllerReference(tr, svc, r.Scheme)
	if err != nil {
		logger.Error(err, "error setting controller reference", "namespace", req.Namespace, "name", serviceName)
		return ctrl.Result{}, err
	}

	logger.Info("creating service", "svc", serviceName)

	err = r.Create(ctx, svc)

	if err != nil {
		logger.Error(err, "failed to create service", "svc", serviceName)
		return ctrl.Result{}, err
	}

	logger.Info("service created successfully", "svc", svc)

	// return response for created successfully
	return ctrl.Result{Requeue: true}, nil
}

func (r *TraefikRedirectReconciler) createIngress(ctx context.Context, req ctrl.Request, tr *traefikv1.TraefikRedirect, ingressName string, serviceName string, traefikType traefikv1.TraefikType, host string) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	pathTypePrefix := networkv1.PathTypePrefix
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
								Path:     "/",
								PathType: &pathTypePrefix,
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
