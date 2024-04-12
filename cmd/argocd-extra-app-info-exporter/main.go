package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	argo "github.com/argoproj/argo-cd/v2/pkg/client/clientset/versioned"
	"github.com/mikejoh/argocd-extra-app-info-exporter/internal/buildinfo"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type exporterOptions struct {
	version    bool
	kubeconfig *string
}

func main() {
	var exporterOpts exporterOptions
	flag.BoolVar(&exporterOpts.version, "version", false, "Print the version number.")
	flag.Parse()

	if exporterOpts.version {
		fmt.Println(buildinfo.Get())
		os.Exit(0)
	}

	if home := homedir.HomeDir(); home != "" {
		exporterOpts.kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		exporterOpts.kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *exporterOpts.kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := argo.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	apps, err := clientset.ArgoprojV1alpha1().Applications("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	for _, app := range apps.Items {
		fmt.Println(app.Name)
	}
}
