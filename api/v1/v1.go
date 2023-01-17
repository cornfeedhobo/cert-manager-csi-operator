/*
Copyright 2023.

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

package v1

import (
	cm "github.com/cornfeedhobo/cert-manager-csi-operator/api/cert-manager"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func mutate(logr logr.Logger, config *cm.Config, objectMeta *metav1.ObjectMeta, podSpec *corev1.PodSpec) {
	if !config.IsManaged(objectMeta) {
		logr.Info("object is not managed")
		return
	}
	logr.Info("object is managed")

	volume, mount := config.GetVolumeAndMount(config.GetAttributes())
	podSpec.Volumes = deepMergeVolumes(podSpec.Volumes, volume)
	for i := range podSpec.Containers {
		podSpec.Containers[i].VolumeMounts = deepMergeMounts(podSpec.Containers[i].VolumeMounts, mount)
	}
	logr.Info("finished ensuring volume and mount exist")
}

func deepMergeVolumes(dst []corev1.Volume, elements ...corev1.Volume) []corev1.Volume {
loop:
	for _, new := range elements {
		for i, current := range dst {
			if current.Name == new.Name {
				dst[i] = new
				continue loop
			}
		}
		dst = append(dst, new)
	}
	return dst
}

func deepMergeMounts(dst []corev1.VolumeMount, elements ...corev1.VolumeMount) []corev1.VolumeMount {
loop:
	for _, new := range elements {
		for i, current := range dst {
			if current.Name == new.Name {
				dst[i] = new
				continue loop
			}
		}
		dst = append(dst, new)
	}
	return dst
}
