package main

import (
	"context"

	endpoints "github.com/deislabs/osiris/pkg/endpoints/hijacker"
	"github.com/deislabs/osiris/pkg/version"
	"k8s.io/klog"
)

func runEndpointsHijacker(ctx context.Context) {
	klog.Infof(
		"Starting Osiris Endpoints Hijacker -- version %s -- commit %s",
		version.Version(),
		version.Commit(),
	)

	cfg, err := endpoints.GetConfigFromEnvironment()
	if err != nil {
		klog.Fatalf(
			"Error retrieving proxy endpoints hijacker webhook server "+
				"configuration: %s",
			err,
		)
	}

	// Run the server
	endpoints.NewHijacker(cfg).Run(ctx)
}
