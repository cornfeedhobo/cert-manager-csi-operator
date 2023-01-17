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
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Deployment", func() {
	It("maybe", func() {
		nn := types.NamespacedName{Namespace: "unittest", Name: "deployment"}

		orig := &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      nn.Name,
				Namespace: nn.Namespace,
				Labels:    map[string]string{"app": "unittest"},
			},
			Spec: appsv1.DeploymentSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app": "unittest"},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "unittest"},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "main",
								Image: "example.com/unittest:latest",
							},
						},
					},
				},
			},
		}
		err := k8sClient.Create(context.TODO(), orig)
		Expect(err).NotTo(HaveOccurred())

		new := &appsv1.Deployment{}
		Eventually(func(g Gomega) {
			err := k8sClient.Get(context.TODO(), nn, new)
			Expect(err).NotTo(HaveOccurred())

			Expect(new.Spec.Template.Spec.Volumes).ToNot(BeNil())
			Expect(new.Spec.Template.Spec.Volumes).To(Equal([]corev1.Volume{volume}))

			Expect(new.Spec.Template.Spec.Containers[0].VolumeMounts).ToNot(BeNil())
			Expect(new.Spec.Template.Spec.Containers[0].VolumeMounts).To(Equal([]corev1.VolumeMount{mount}))
		}, time.Second*10, time.Millisecond*250).Should(Succeed())
	})
})
