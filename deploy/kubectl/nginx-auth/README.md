## Ingress-Nginx Basic Auth

This can be used in the cluster which uses the nginx-ingress.

We use `basic-auth` to control the access of kube-explorer. The auth token is stored in the secret.

The default user is `niusmallnan`, and password is `dagedddd`. You can replace to another value with `htpasswd` tool.

```
htpasswd -nb username password | base64
```

To install this mode, just run this script:

```
kubectl create -f ./secret.yaml
export MY_IP=$(curl -sL ipinfo.io/ip)
envsubst < ./ingress.yaml.tpl | kubectl create -f -
```

For more infos: https://kubernetes.github.io/ingress-nginx/examples/auth/basic/
