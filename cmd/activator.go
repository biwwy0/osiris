package main

import (
	"context"

	deployments "github.com/deislabs/osiris/pkg/deployments/activator"
	"github.com/deislabs/osiris/pkg/kubernetes"
	"github.com/deislabs/osiris/pkg/version"
	"k8s.io/klog"
)

func runActivator(ctx context.Context) {
	klog.Infof(
		"Starting Osiris Activator -- version %s -- commit %s",
		version.Version(),
		version.Commit(),
	)

	client, err := kubernetes.Client()
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err)
	}

	activator, err := deployments.NewActivator(client)
	if err != nil {
		klog.Fatalf("Error initializing activator: %s", err)
	}

	// Run the activator
	activator.Run(ctx)
}
