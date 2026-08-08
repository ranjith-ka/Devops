#!/usr/bin/env sh
set -eu

: "${IMAGE:?Skaffold must provide IMAGE}"
: "${BUILD_CONTEXT:?Skaffold must provide BUILD_CONTEXT}"
: "${KUBECONTEXT:?Skaffold must provide KUBECONTEXT}"
: "${NAMESPACE:?Skaffold must provide NAMESPACE}"
: "${DOCKER_CONFIG_SECRET_NAME:?Skaffold must provide DOCKER_CONFIG_SECRET_NAME}"

BUILDKIT_IMAGE="${BUILDKIT_IMAGE:-moby/buildkit:v0.32.2-rootless}"
CACHE_IMAGE="${CACHE_IMAGE:-ranjithka/dotnet-skaffold-api:buildcache}"
DOCKER_CONFIG_PATH="${DOCKER_CONFIG_PATH:-$HOME/.docker/config.json}"
APP_NAMESPACE="${APP_NAMESPACE:-dotnet-skaffold}"
POD_NAME="buildkit-$(date +%s)-$$"

if [ ! -f "$DOCKER_CONFIG_PATH" ]; then
  echo "Docker registry config not found: $DOCKER_CONFIG_PATH" >&2
  echo "Run 'docker login' before starting Skaffold." >&2
  exit 1
fi

cleanup() {
  status=$?
  if [ "$status" -ne 0 ] && [ "${KEEP_FAILED_BUILDER:-true}" = "true" ]; then
    echo "Build failed; preserving pod $POD_NAME for diagnostics" >&2
    echo "Delete it with: kubectl --context $KUBECONTEXT -n $NAMESPACE delete pod $POD_NAME" >&2
    return
  fi
  kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" \
    delete pod "$POD_NAME" --ignore-not-found --wait=false >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

echo "Starting rootless BuildKit pod $POD_NAME in $KUBECONTEXT/$NAMESPACE"

# Recreate prerequisites when the disposable Kind cluster has been recreated.
kubectl --context "$KUBECONTEXT" create namespace "$NAMESPACE" \
  --dry-run=client -o yaml | \
  kubectl --context "$KUBECONTEXT" apply -f -

kubectl --context "$KUBECONTEXT" create namespace "$APP_NAMESPACE" \
  --dry-run=client -o yaml | \
  kubectl --context "$KUBECONTEXT" apply -f -

kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" \
  create secret generic "$DOCKER_CONFIG_SECRET_NAME" \
  --from-file="config.json=$DOCKER_CONFIG_PATH" \
  --dry-run=client -o yaml | \
  kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" apply -f -

# Keep BuildKit's content store between ephemeral builder pods. The cluster's
# default StorageClass provisions this claim on first use.
kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: buildkit-cache
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 10Gi
EOF

kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" create -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  name: ${POD_NAME}
  labels:
    app.kubernetes.io/name: skaffold-buildkit
spec:
  restartPolicy: Never
  securityContext:
    fsGroup: 1000
  containers:
    - name: buildkit
      image: ${BUILDKIT_IMAGE}
      imagePullPolicy: IfNotPresent
      command: ["sh", "-c", "touch /tmp/ready && sleep 86400"]
      env:
        - name: DOCKER_CONFIG
          value: /home/user/.docker
        - name: BUILDKITD_FLAGS
          value: --oci-worker-no-process-sandbox
      readinessProbe:
        exec:
          command: ["test", "-f", "/tmp/ready"]
        initialDelaySeconds: 1
        periodSeconds: 1
      securityContext:
        runAsUser: 1000
        runAsGroup: 1000
        seccompProfile:
          type: Unconfined
        appArmorProfile:
          type: Unconfined
      volumeMounts:
        - name: workspace
          mountPath: /workspace
        - name: buildkit-state
          mountPath: /home/user/.local/share/buildkit
        - name: registry-auth
          mountPath: /home/user/.docker/config.json
          subPath: config.json
          readOnly: true
  volumes:
    - name: workspace
      emptyDir: {}
    - name: buildkit-state
      persistentVolumeClaim:
        claimName: buildkit-cache
    - name: registry-auth
      secret:
        secretName: ${DOCKER_CONFIG_SECRET_NAME}
EOF

kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" \
  wait --for=condition=Ready "pod/$POD_NAME" --timeout="${TIMEOUT:-300s}"

echo "Uploading build context to $POD_NAME"
tar \
  --exclude=.git \
  --exclude=bin \
  --exclude=obj \
  -C "$BUILD_CONTEXT" -cf - . | \
  kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" \
    exec -i "$POD_NAME" -- tar -xf - -C /workspace

echo "Building and pushing $IMAGE from inside Kubernetes"
kubectl --context "$KUBECONTEXT" --namespace "$NAMESPACE" \
  exec "$POD_NAME" -- buildctl-daemonless.sh build \
    --progress=plain \
    --frontend=dockerfile.v0 \
    --local=context=/workspace \
    --local=dockerfile=/workspace \
    --import-cache="type=registry,ref=$CACHE_IMAGE" \
    --output="type=image,name=$IMAGE,push=true" \
    --export-cache="type=registry,ref=$CACHE_IMAGE,mode=max,ignore-error=true"

echo "Remote BuildKit build completed: $IMAGE"
