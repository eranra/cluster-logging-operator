package runtime

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceBuilder struct {
	obj *corev1.Service
}

func NewServiceBuilder(svc *corev1.Service) *ServiceBuilder {
	return &ServiceBuilder{
		obj: svc,
	}
}

func (builder *ServiceBuilder) WithSelector(selector map[string]string) *ServiceBuilder {
	builder.obj.Spec.Selector = selector
	return builder
}

func (builder *ServiceBuilder) AddServicePort(port int32, targetPort int) *ServiceBuilder {
	builder.obj.Spec.Ports = append(builder.obj.Spec.Ports, corev1.ServicePort{
		Port:       port,
		TargetPort: intstr.FromInt(targetPort),
	})
	return builder
}
