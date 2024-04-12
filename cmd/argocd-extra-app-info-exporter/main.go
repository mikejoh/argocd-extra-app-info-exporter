package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikejoh/argocd-extra-app-info-exporter/internal/buildinfo"
)

type argocd-extra-app-info-exporterOptions struct {
	version bool
}

func main() {
	var argocd-extra-app-info-exporterOpts argocd-extra-app-info-exporterOptions
	flag.BoolVar(&argocd-extra-app-info-exporterOpts.version, "version", false, "Print the version number.")
	flag.Parse()

	if argocd-extra-app-info-exporterOpts.version {
		fmt.Println(buildinfo.Get())
		os.Exit(0)
	}
}
