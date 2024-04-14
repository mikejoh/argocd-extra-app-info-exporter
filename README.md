# argocd-extra-app-info-exporter

`argocd-extra-app-info-exporter` - Exports that one missing metric from ArgoCD.

This exporter exports a `argocd_extra_app_info` metric which looks like the original `argocd_app_info` metric but with the `targetRevision` field as the only label.

One of the reasons one would want that label is to create e.g. alerts if the `targetRevision` field is anything else than `main`, this exporter will probably be around until https://github.com/argoproj/argo-cd/pull/15143 merges!

## Install

```
helm repo add mikejoh https://mikejoh.github.io/helm-charts/

helm upgrade \
  --install
  --namespace argocd-extra-app-info-exporter \
  --create-namespace \
  mikejoh/argocd-extra-app-info-exporter \
  argocd-extra-app-info-exporter
```

Add `--set serviceMonitor.enabled=true` to deploy a `ServiceMonitor` (part of the Prometheus Operator).
