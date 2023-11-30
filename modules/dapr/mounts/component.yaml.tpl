apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: {{ .Name }}
spec:
  type: {{ .Type }}
  version: v1
  {{- if .Metadata }}
  metadata:
    {{- range $key, $value := .Metadata }}
    - name: {{ $key }}
      value: {{ $value }}
    {{- end }}
  {{- end }}
