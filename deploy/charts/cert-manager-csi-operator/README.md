# cert-manager-csi-operator

This chart is only for demonstration purposes.

## Usage

To try out this chart, you'll need to first install the cert-manager service:

```shell
helm repo add jetstack https://charts.jetstack.io
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --set installCRDs=true
```

Then you can install this chart:

```shell
helm install "$(basename "${PWD}")" .
```

To test, create a Deployment you can inspect:

```shell
kubectl apply -f <(cat <<'EOF'
apiVersion: "apps/v1"
kind: "Deployment"
metadata:
  name: "example"
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: "example"
  template:
    metadata:
      labels:
        app.kubernetes.io/name: "example"
    spec:
      containers:
        - name: "main"
          image: "alpine:latest"
          args: ["sleep", "600"]
EOF
)
```

Then check out the created resource:

```shell
kubectl get deployments/example --output=jsonpath='{.metadata.annotations}' | jq .
```
