package runtime

import (
	corev1 "k8s.io/api/core/v1"
)

type ConfigMapBuilder struct {
	obj *corev1.ConfigMap
}

func NewConfigMapBuilder(cm *corev1.ConfigMap) *ConfigMapBuilder {
	return &ConfigMapBuilder{
		obj: cm,
	}
}
func (builder *ConfigMapBuilder) Add(key, value string) *ConfigMapBuilder {
	builder.obj.Data[key] = value
	return builder
}
