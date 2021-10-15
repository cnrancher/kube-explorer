# kube-explorer

kube-explorer is a portable explorer for Kubernetes without any dependency.

It integrates the Rancher steve framework and its dashboard, and is recompiled, packaged, compressed, and provides an almost completely stateless Kubernetes resource manager.

## Usage ✅

Please download the binary from the [release page](https://github.com/niusmallnan/kube-explorer/releases).

To run an HTTP only server:

```
./kube-explorer --kubeconfig=xxxx --http-listen-port=9898 --https-listen-port=0
```

Then, open the browser to visit http://x.x.x.x:9898 .

## Build ✅

To debug on an AMD64 Linux host:

```
make dev

# $basedir=/opt/ui/dist/
# prepare the file trees like this
# $basedir/dashboard/
# $basedir/index.html

# good to go!
./kube-explorer --debug  --ui-path /opt/ui/dist/ --http-listen-port=9898 --https-listen-port=0
```

To build all cross-platform binaries:

```
CROSS=1 make
```
