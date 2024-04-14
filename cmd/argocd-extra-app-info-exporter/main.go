package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	argo "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"github.com/mikejoh/argocd-extra-app-info-exporter/internal/buildinfo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type exporterOptions struct {
	version              bool
	interval             DurationFlag
	metricsListenAddress string
	metricsPath          string
}

const (
	// DurationFlag is the default interval for the exporter.
	defaultDuration = DurationFlag(1 * time.Minute)
)

func main() {
	var exporterOpts exporterOptions
	flag.BoolVar(&exporterOpts.version, "version", false, "Print the version number.")
	flag.StringVar(&exporterOpts.metricsListenAddress, "metrics-listen-address", "0.0.0.0:9999", "Set the metrics listen address.")
	flag.StringVar(&exporterOpts.metricsPath, "metrics-path", "/metrics", "Set the metrics path.")
	flag.Var(&exporterOpts.interval, "interval", "Application fetch interval in human-friendly format (e.g., 5s for 5 seconds, 10m for 10 minutes)")
	flag.Parse()

	if exporterOpts.version {
		fmt.Println(buildinfo.Get())
		os.Exit(0)
	}

	if exporterOpts.interval == 0 {
		exporterOpts.interval = defaultDuration
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := argo.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	appExtraInfo := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "argocd_extra_app_info",
		Help: "Extra information about application.",
	}, []string{
		"namespace",
		"name",
		"project",
		"targetRevision",
	})

	prometheus.MustRegister(appExtraInfo)

	log.Printf("starting argocd-extra-app-info-exporter %s, fetching application(s) every %s", buildinfo.Get(), exporterOpts.interval.String())
	go func() {
		ticker := time.NewTicker(time.Duration(exporterOpts.interval))
		for {
			select {
			case <-ticker.C:
				apps, err := clientset.ArgoprojV1alpha1().Applications("").List(context.Background(), v1.ListOptions{})
				if err != nil {
					log.Fatal(err)
				}
				for _, app := range apps.Items {
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
		log.Fatalf("listen: %s\n", err)
	}
}
