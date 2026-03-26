# kubectl CLI image — installs the official release binary from dl.k8s.io (same bits the
# Kubernetes project publishes from kubernetes/kubernetes). The kubernetes/kubectl repo is
# a mirror/issue tracker; see https://github.com/kubernetes/kubectl/blob/master/README.md
#
# Build:
#   docker build -t my-kubectl:v1.31.3 --build-arg KUBECTL_VERSION=v1.31.3 .
#
# Multi-arch (push to registry):
#   docker buildx build --platform linux/amd64,linux/arm64 -t my-kubectl:v1.31.3 --push .

FROM alpine:3.20

ARG KUBECTL_VERSION=v1.31.3

RUN apk add --no-cache ca-certificates curl \
    && ARCH=$(uname -m) \
    && case "$ARCH" in \
         aarch64|arm64) KARCH=arm64 ;; \
         x86_64|amd64) KARCH=amd64 ;; \
         *) echo "unsupported arch: $ARCH" >&2; exit 1 ;; \
       esac \
    && curl -fsSL -o /usr/local/bin/kubectl \
         "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/${KARCH}/kubectl" \
    && chmod +x /usr/local/bin/kubectl \
    && kubectl version --client=true \
    && rm -rf /var/cache/apk/*

ENTRYPOINT ["/usr/local/bin/kubectl"]
