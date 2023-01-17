{{/*
	Expand the name.
*/}}
{{- define "common.name" -}}
{{- tpl .Values.name . -}}
{{- end -}}

{{/*
	Create chart name and version as used by the chart label.
*/}}
{{- define "common.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
	Selector labels
*/}}
{{- define "common.selectorLabels" -}}
app.kubernetes.io/name: {{ include "common.name" . | quote }}
app.kubernetes.io/instance: {{ .Release.Name | quote }}
{{- end -}}

{{/*
	Common labels
*/}}
{{- define "common.labels" -}}
helm.sh/chart: {{ include "common.chart" . | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service | quote }}
{{ include "common.selectorLabels" . }}
{{- end -}}

{{/*
	Pod labels
*/}}
{{- define "common.versionLabel" -}}
{{- if .Chart.AppVersion -}}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end -}}
{{- end -}}
