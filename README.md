# argocd-extra-app-info-exporter

[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/mikejoh)](https://artifacthub.io/packages/search?repo=mikejoh)

`argocd-extra-app-info-exporter` - Exports that one missing metric from ArgoCD.

This exporter exports a `argocd_extra_app_info` metric which looks like the original `argocd_app_info` metric but with the application `revision` field as a label. One of the reasons one would want that label is to create e.g. alerts if the `revision` field is anything else than `main`, this exporter will probably be around until PR [`15143`](https://github.com/argoproj/argo-cd/pull/15143) merges.

Exported metric labels:
* `namespace`
* `name`
* `project`
* `revision`

_Please note that this exporter lists all application in a cluster (once per interval), if you've specified a namespace the list of applications will be limited to that namespace._

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

### Flags

* Add `--set serviceMonitor.enabled=true` to deploy a `ServiceMonitor` (part of the Prometheus Operator).
* Add `--set prometheusRule.enabled=true` to deploy a proof-of-concept [alert rule](https://github.com/mikejoh/helm-charts/blob/main/charts/argocd-extra-app-info-exporter/templates/prometheusrule.yaml#L9-L25) (`PrometheusRule`, also part of the Prometheus Operator).
* Add `--set excludeRevisions[0]="main"` to exclude creating a metric for the revision `main`.
