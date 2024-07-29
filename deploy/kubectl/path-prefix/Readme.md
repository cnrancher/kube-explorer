# Deploy kube-explorer behind proxy with path prefix

> Supported since v0.5.0

The kube-explorer dashboard can be exposed behind a proxy and path prefix like `http://your-domain.com/kube-explorer`.

The deployment examples in this folder are:

- `nginx ingress`
- `traefik ingress`

## Serve with ingress

When serving with nginx/traefik ingress controller, the template ingress file needs to be modified. In the `*.tpl` file, you can spot the missing hostname like:

```yaml
spec:
  rules:
    - host: "${MY_IP}.sslip.io" # Replace with your actual domain
```

Replace your ip to `${MY_UP}`, this will use the [sslip.io](https://sslip.io/) dns service to resolve the hostname to the ingress ip address.

For the traefik ingress, it is using `v2` version of the traefik ingress schema which use middlewares to modify the proxy request. Both `stripPrefix` and `headers` are used.  
For the nginx ingress, the annotations `nginx.ingress.kubernetes.io/x-forwarded-prefix` and `nginx.ingress.kubernetes.io/rewrite-target` are used to strip prefix and to add proxy request header.

## Serve with self-hosted proxy

If serving the kube-explorer with self-hosted proxy, following modifications are required when proxying:

- Rewrite the proxy request to strip the path prefix like `rewrite "(?i)/kube-explorer(/|$)(.*)" /$2 break;` in nginx configuration.
- Add header `X-API-URL-Prefix` or `X-Forwarded-Prefix` with the path prefix when proxying request like `proxy_set_header X-Forwarded-Prefix "/kube-explorer";` in nginx configuration.

Then kube-explorer will response the index.html with modified content with path prefix to the browser.
