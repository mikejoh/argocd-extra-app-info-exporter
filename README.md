# argocd-extra-app-info-exporter

`argocd-extra-app-info-exporter` - Exports that one missing metric from ArgoCD.

This exporter exports a `argocd_extra_app_info` metric which looks like the already exported `argocd_app_info` metric but with the `targetRevision` field as label.

This exporter will probably be around until https://github.com/argoproj/argo-cd/pull/15143 merges!

## Install

1. `make build`
