package ingress

import (
	"slices"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

type Config struct {
	Name                string
	Namespace           string
	Labels              map[string]string
	IngressClassName    string
	DefaultBackend      *networkingv1.IngressBackend
	HostName            string
	AdditionalHostNames []string
	Service             string
}

func New(cfg Config) *networkingv1.Ingress {
	return &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: networkingv1.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Name,
			Namespace: cfg.Namespace,
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: ptr.To(cfg.IngressClassName),
			DefaultBackend:   cfg.DefaultBackend,
			TLS: []networkingv1.IngressTLS{
				{
					Hosts: slices.Concat([]string{cfg.HostName}, cfg.AdditionalHostNames),
				},
			},
		},
	}
}
