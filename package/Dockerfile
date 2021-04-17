FROM alpine:3.13

COPY kube-explorer entrypoint.sh /usr/bin/
# Hack to make golang do files,dns search order
ENV LOCALDOMAIN=""
ENTRYPOINT ["entrypoint.sh"]
