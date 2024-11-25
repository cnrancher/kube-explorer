FROM registry.suse.com/bci/bci-base:15.6
ARG TARGETARCH
ARG TARGETOS
ENV ARCH=${TARGETARCH:-"amd64"} OS=${TARGETOS:-"linux"}
RUN zypper install -y catatonit
COPY kube-explorer-${OS}-${ARCH} /usr/bin/kube-explorer
ENTRYPOINT [ "/usr/bin/catatonit", "--", "kube-explorer" ]
