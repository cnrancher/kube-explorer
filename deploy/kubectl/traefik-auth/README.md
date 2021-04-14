## Traefik Auth

This can be used in K3s, as K3s use traefik as the default ingress class.

We use `basic-auth` to control the access of kube-explorer. The auth token is stored in the secret.

The default user is `niusmallnan`, and password is `dagedddd`. You can replace to another value with `htpasswd` tool.

```
htpasswd -nb username password | base64
```

To install this mode, just run this script:

```
kubectl apply -f ./secret.yaml
export MY_XIP_IO=$(curl -sL ipinfo.io/ip)
envsubst < ./ingress.yaml.tpl | kubectl apply -f -
```

For more infos: https://doc.traefik.io/traefik/v1.7/configuration/backends/kubernetes/
