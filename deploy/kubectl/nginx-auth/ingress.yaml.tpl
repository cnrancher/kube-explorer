# Note: please replace the host first
# To use xip.io: http://xip.io/
# To get your public IP: curl ipinfo.io/ip

apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: kube-explorer
  namespace: kube-system
  labels:
    app: kube-explorer
  annotations:
    kubernetes.io/ingress.class: "nginx"
    nginx.ingress.kubernetes.io/auth-type: basic
    nginx.ingress.kubernetes.io/auth-secret: kube-explorer
    nginx.ingress.kubernetes.io/auth-realm: 'Authentication Required - kube-explorer'
spec:
  rules:
  - host: "${MY_XIP_IO}.xip.io"
    http:
      paths:
      - path: /
        backend:
          serviceName: kube-explorer
          servicePort: 8989
