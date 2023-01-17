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
	"encoding/json"
	"net/http"

	cm "github.com/cornfeedhobo/cert-manager-csi-operator/api/cert-manager"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/mutate-v1-statefulset,mutating=true,failurePolicy=fail,sideEffects=None,groups=apps,resources=statefulsets,verbs=create;update,versions=v1,name=mstatefulset.kb.io,admissionReviewVersions=v1

type StatefulSetMutator struct {
	Client  client.Client
	Config  *cm.Config
	Logger  logr.Logger
	decoder *admission.Decoder
}

func RegisterStatefulSetMutatorWebhook(mgr manager.Manager, c *cm.Config) {
	l := logf.Log.WithName("statefulset-resource")
	l.Info("Registering StatefulSetMutator")

	m := StatefulSetMutator{
		Client: mgr.GetClient(),
		Config: c,
		Logger: l,
	}
	mgr.GetWebhookServer().Register("/mutate-v1-statefulset", &webhook.Admission{Handler: &m})
}

func (m *StatefulSetMutator) InjectDecoder(d *admission.Decoder) error {
	m.decoder = d
	return nil
}

func (m *StatefulSetMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	obj := &appsv1.StatefulSet{}

	err := m.decoder.Decode(req, obj)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	mutate(m.Logger, m.Config, &obj.ObjectMeta, &obj.Spec.Template.Spec)

	marshaled, err := json.Marshal(obj)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaled)
}
