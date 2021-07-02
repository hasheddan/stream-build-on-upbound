# Build on Upbound

## Setup

1. Install ArgoCD into Cluster

```
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

2. Install ArgoCD CLI

```
curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
chmod +x /usr/local/bin/argocd
```

3. Install Upbound CLI

```
curl -sL https://cli.upbound.io | sh
```

4. Install Upbound Universal Crossplane (UXP) into Cluster

```
up uxp install
```

5. Attach UXP to Upbound Cloud

```
up cloud ctp attach my-ctp | up uxp connect -
```

5. Sync App of Apps

```
argocd app create apps \
    --dest-namespace argocd \
    --dest-server https://kubernetes.default.svc \
    --repo https://github.com/hasheddan/stream-build-on-upbound.git \
    --path apps/local  
argocd app sync apps  
```
