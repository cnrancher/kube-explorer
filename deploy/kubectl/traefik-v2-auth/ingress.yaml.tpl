# Note: please replace the host first
# To use sslip.io.io: https://sslip.io.io/
# To get your public IP: curl ipinfo.io/ip

apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: kube-explorer
  namespace: kube-system
  labels:
    app: kube-explorer
  annotations:
    traefik.ingress.kubernetes.io/router.middlewares: kube-system-kube-explorer@kubernetescrd
spec:
  rules:
  - host: "${MY_IP}.sslip.io"
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: kube-explorer
            port:
              number: 8989
