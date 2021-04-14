FROM golang:1.15.11 as build
ENV GIT_COMMIT="255623bd7c0277f8b8c61e66fdf873d5a43778a3"
ENV GIT_BRANCH="kube-explorer"
RUN \
    mkdir -p /go/src/github.com/rancher/ && \
    cd /go/src/github.com/rancher/ && \
    git clone --depth=1 --branch ${GIT_BRANCH} https://github.com/niusmallnan/steve.git && \
    cd steve && \
    git reset --hard ${GIT_COMMIT} && \
    go mod vendor && \
    CGO_ENABLED=0 go build -ldflags "-X github.com/rancher/steve/pkg/ui.UIOffline=true -extldflags -static -s" -o /kube-explorer

FROM alpine:3.13

ENV CATTLE_DASHBOARD_UI_VERSION v2.5.8-rc3

RUN apk -U --no-cache add bash curl ca-certificates && \
    mkdir -p /usr/share/rancher/ui-dashboard/dashboard && \
    cd /usr/share/rancher/ui-dashboard/dashboard && \
    curl -sL https://releases.rancher.com/dashboard/${CATTLE_DASHBOARD_UI_VERSION}.tar.gz | tar xvzf - --strip-components=2 && \
    ln -s dashboard/index.html ../index.html

COPY --from=build /kube-explorer /usr/bin/
COPY entrypoint.sh /usr/bin
# Hack to make golang do files,dns search order
ENV LOCALDOMAIN=""
ENTRYPOINT ["entrypoint.sh"]
