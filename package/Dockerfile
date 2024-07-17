FROM registry.suse.com/bci/bci-minimal:15.6
ARG TARGETARCH
ARG TARGETOS
ENV ARCH=${TARGETARCH:-"amd64"} OS=${TARGETOS:-"linux"}
COPY entrypoint.sh /usr/bin/
COPY kube-explorer-${OS}-${ARCH} /usr/bin/kube-explorer
ENTRYPOINT ["entrypoint.sh"]
