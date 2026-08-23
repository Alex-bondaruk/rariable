{{- define "rariable.fullname" -}}
{{ .Release.Name }}-{{ .Chart.Name }}
{{- end -}}

{{- define "rariable.chart" -}}
{{ .Chart.Name }}-{{ .Chart.Version }}
{{- end -}}

{{/* Identity labels only — used wherever Kubernetes matches pods (immutable). */}}
{{- define "rariable.selectorLabels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/* Full metadata labels — informational, safe to change across upgrades. */}}
{{- define "rariable.labels" -}}
{{ include "rariable.selectorLabels" . }}
helm.sh/chart: {{ include "rariable.chart" . }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "rariable.secretName" -}}
{{- if .Values.secret.existingSecret -}}
{{ .Values.secret.existingSecret }}
{{- else -}}
{{ include "rariable.fullname" . }}
{{- end -}}
{{- end -}}
