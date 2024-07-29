# Note: please replace the host first
# To use sslip.io: https://sslip.io/
# To get your public IP: curl ipinfo.io/ip
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/x-forwarded-prefix: "/kube-explorer"
    nginx.ingress.kubernetes.io/rewrite-target: /$2
  name: kube-explorer-ingress
  namespace: kube-system
spec:
  rules:
    - host: "${MY_IP}.sslip.io" # Replace with your actual domain
      http:
        paths:
          - backend:
              service:
                name: kube-explorer
                port:
                  name: http
            path: /kube-explorer(/|$)(.*)
            pathType: ImplementationSpecific

