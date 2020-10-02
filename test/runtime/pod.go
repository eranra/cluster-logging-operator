package runtime

import (
	corev1 "k8s.io/api/core/v1"
)

type PodBuilder struct {
	obj *corev1.Pod
}

func NewPodBuilder(pod *corev1.Pod) *PodBuilder {

	return &PodBuilder{
		obj: pod,
	}
}

type ContainerBuilder struct {
	container  corev1.Container
	podBuilder *PodBuilder
}

func (builder *ContainerBuilder) End() *PodBuilder {
	builder.podBuilder.obj.Spec.Containers = append(builder.podBuilder.obj.Spec.Containers, builder.container)
	return builder.podBuilder
}

func (builder *ContainerBuilder) AddVolumeMount(name, path, subPath string, readonly bool) *ContainerBuilder {
	builder.container.VolumeMounts = append(builder.container.VolumeMounts, corev1.VolumeMount{
		Name:      name,
		ReadOnly:  readonly,
		MountPath: path,
		SubPath:   subPath,
	})
	return builder
}

func (builder *ContainerBuilder) AddEnvVar(name, value string) *ContainerBuilder {
	builder.container.Env = append(builder.container.Env, corev1.EnvVar{
		Name:  name,
		Value: value,
	})
	return builder
}
func (builder *ContainerBuilder) AddEnvVarFromFieldRef(name, fieldRef string) *ContainerBuilder {
	builder.container.Env = append(builder.container.Env, corev1.EnvVar{
		Name: name,
		ValueFrom: &corev1.EnvVarSource{
			FieldRef: &corev1.ObjectFieldSelector{
				FieldPath: fieldRef,
			},
		},
	})
	return builder
}

func (builder *PodBuilder) AddContainer(name, image string) *ContainerBuilder {
	containerBuilder := ContainerBuilder{
		container: corev1.Container{
			Name:  name,
			Image: image,
			Env:   []corev1.EnvVar{},
		},
		podBuilder: builder,
	}
	return &containerBuilder
}

func (builder *PodBuilder) AddConfigMapVolume(name, configMapName string) *PodBuilder {
	builder.obj.Spec.Volumes = append(builder.obj.Spec.Volumes, corev1.Volume{
		Name: name,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: configMapName,
				},
			},
		},
	})
	return builder
}

func (builder *PodBuilder) WithLabels(labels map[string]string) *PodBuilder {
	builder.obj.Labels = labels
	return builder
}