FROM registry.suse.com/bci/bci-minimal:15.5
ARG TARGETARCH
ARG TARGETOS
ENV ARCH=${TARGETARCH:-"amd64"} OS=${TARGETOS:-"linux"}
COPY package/entrypoint.sh /usr/bin/
COPY bin/kube-explorer-${OS}-${ARCH} /usr/bin/kube-explorer
ENTRYPOINT ["entrypoint.sh"]
