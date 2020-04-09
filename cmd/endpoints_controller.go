package main

import (
	"context"

	endpoints "github.com/deislabs/osiris/pkg/endpoints/controller"
	"github.com/deislabs/osiris/pkg/kubernetes"
	"github.com/deislabs/osiris/pkg/version"
	"k8s.io/klog"
)

func runEndpointsController(ctx context.Context) {
	klog.Infof(
		"Starting Osiris Endpoints Controller -- version %s -- commit %s",
		version.Version(),
		version.Commit(),
	)

	client, err := kubernetes.Client()
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err)
	}

	controllerCfg, err := endpoints.GetConfigFromEnvironment()
	if err != nil {
		klog.Fatalf(
			"Error retrieving endpoints controller configuration: %s",
			err,
		)
	}

	// Run the controller
	endpoints.NewController(controllerCfg, client).Run(ctx)
}
