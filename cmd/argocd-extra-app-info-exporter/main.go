package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"log/slog"

	argo "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"github.com/mikejoh/argocd-extra-app-info-exporter/internal/buildinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type exporterOptions struct {
	version              bool
	interval             DurationFlag
	metricsListenAddress string
	metricsPath          string
	namespace            string
	excludeRevisions     string
}

const (
	defaultDuration = DurationFlag(1 * time.Minute)
)

func main() {
	var exporterOpts exporterOptions
	flag.BoolVar(&exporterOpts.version, "version", false, "Print the version number.")
	flag.StringVar(&exporterOpts.metricsListenAddress, "metrics-listen-address", "0.0.0.0:9999", "Set the metrics listen address.")
	flag.StringVar(&exporterOpts.metricsPath, "metrics-path", "/metrics", "Set the metrics path.")
	flag.StringVar(&exporterOpts.namespace, "namespace", "", "List all applications from this namespace. Default is all namespaces.")
	flag.StringVar(&exporterOpts.excludeRevisions, "exclude-revisions", "", "Comma-separated list of revisions to exclude from the metrics. Example: 'main,master,HEAD'")
	flag.Var(&exporterOpts.interval, "interval", "Application fetch interval in human-friendly format (e.g., 5s for 5 seconds, 10m for 10 minutes)")
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if exporterOpts.version {
		fmt.Println(buildinfo.Get())
		os.Exit(0)
	}

	var revs []string
	if exporterOpts.excludeRevisions != "" {
		revs = strings.Split(exporterOpts.excludeRevisions, ",")
		logger.Info("excluding the following revisions from metrics export", "revisions", exporterOpts.excludeRevisions)
	}

	if exporterOpts.interval == 0 {
		exporterOpts.interval = defaultDuration
	}

	clientset, err := getClientset()
	if err != nil {
		logger.Error("failed to create clientset", "err", err)
		os.Exit(1)
	}

	appExtraInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "argocd_extra_app_info",
		Help: "Extra information about application.",
	}, []string{
		"namespace",
		"name",
		"project",
		"revision",
	})

	prometheus.MustRegister(appExtraInfo)

	bi := buildinfo.Get()
	logger.Info("starting", "app", buildinfo.Get().Name, "version", bi.Version, "interval", exporterOpts.interval.String())
	go func() {
		ticker := time.NewTicker(time.Duration(exporterOpts.interval))
		for {
			select {
			case <-ticker.C:
				apps, err := clientset.ArgoprojV1alpha1().Applications(exporterOpts.namespace).List(context.Background(), v1.ListOptions{})
				if err != nil {
					logger.Warn("failed to list applications", "err", err)
					continue
				}

				if len(apps.Items) == 0 {
					logger.Info("no applications found", "namespace", exporterOpts.namespace)
					continue
				}

				logger.Info("applications found", "num", len(apps.Items))

				for _, app := range apps.Items {
					// Skip if revision is not set
					if app.Spec.GetSource().TargetRevision == "" {
						continue
					}

					// Skip if revision is in the exclude list
					if slices.Contains(revs, app.Spec.GetSource().TargetRevision) {
						continue
					}

					appExtraInfo.WithLabelValues(
						app.Namespace,
						app.Name,
						app.Spec.GetProject(),
						app.Spec.GetSource().TargetRevision,
					).Set(1)
				}
			}
		}
	}()

	mux := http.NewServeMux()
	mux.Handle(exporterOpts.metricsPath, promhttp.Handler())

	httpServer := &http.Server{
		Addr:         exporterOpts.metricsListenAddress,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("starting HTTP server failed", "err", err)
		os.Exit(1)
	}
}

func getClientset() (*argo.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		var kubeconfig string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}

		return argo.NewForConfig(config)
	}

	return argo.NewForConfig(config)
}
