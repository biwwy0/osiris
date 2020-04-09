package main

import (
	"context"

	deployments "github.com/deislabs/osiris/pkg/deployments/zeroscaler"
	"github.com/deislabs/osiris/pkg/kubernetes"
	"github.com/deislabs/osiris/pkg/version"
	"k8s.io/klog"
)

func runZeroScaler(ctx context.Context) {
	klog.Infof(
		"Starting Osiris Zeroscaler -- version %s -- commit %s",
		version.Version(),
		version.Commit(),
	)

	client, err := kubernetes.Client()
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	cfg, err := deployments.GetConfigFromEnvironment()
	if err != nil {
		klog.Fatalf("Error getting zeroscaler envconfig: %s", err.Error())
	}

	// Run the zeroscaler
	deployments.NewZeroscaler(cfg, client).Run(ctx)
}
