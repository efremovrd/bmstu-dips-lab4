{{- define "deployment.template" }}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .ctx.Release.Name }}-{{.service.name}}-dep
  labels:
    app: {{ .ctx.Release.Name }}-{{.service.name}}
spec:
  replicas: {{.service.replicaCount}}
  selector:
    matchLabels:
      app: {{ .ctx.Release.Name }}-{{.service.name}}
  template:
    metadata:
      name: {{ .ctx.Release.Name }}-{{.service.name}}
      labels:
        app: {{ .ctx.Release.Name }}-{{.service.name}}
    spec:
      containers:
        - name: {{ .ctx.Release.Name }}-{{.service.name}}
          image: {{.service.container}}
          imagePullPolicy: Always
      restartPolicy: Always
{{- end}}

