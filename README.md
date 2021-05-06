# kube-explorer

kube-explorer is portable explorer for Kubernetes, without any dependency.

It integrates the Rancher steve framework and its dashboard, and is recompiled, packaged, compressed, and provides an almost completely stateless Kubernetes resource manager.

## Usage ✅

Please download the binary from the [release page](https://github.com/niusmallnan/kube-explorer/releases).

To run an HTTP only server:

```
./kube-explorer --kubeconfig=xxxx --http-listen-port=9898 --https-listen-port=0
```

Then, open the browser to visit http://x.x.x.x:9898 .
