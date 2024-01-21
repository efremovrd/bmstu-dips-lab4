{{- define "service.template" }}
apiVersion: v1
kind: Service
metadata:
  name: {{.ctx.Release.Name}}-{{.service.name}}-srv
spec:
  selector:
    app: {{.ctx.Release.Name}}-{{.service.name}}
  ports:
    - protocol: TCP
      port: {{.service.port}}
      targetPort: {{.service.port}}
  type: NodePort
{{- end }}

