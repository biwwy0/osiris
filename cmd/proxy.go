package main

import (
	"context"

	"github.com/deislabs/osiris/pkg/metrics/proxy/proxy"
	"github.com/deislabs/osiris/pkg/version"
	"k8s.io/klog"
)

func runProxy(ctx context.Context) {
	klog.Infof(
		"Starting Osiris Proxy -- version %s -- commit %s",
		version.Version(),
		version.Commit(),
	)

	cfg, err := proxy.GetConfigFromEnvironment()
	if err != nil {
		klog.Fatalf("Error retrieving proxy configuration: %s", err)
	}

	proxy, err := proxy.NewProxy(cfg)
	if err != nil {
		klog.Fatalf("Error initializing proxy: %s", err)
	}

	// Run the proxy
	proxy.Run(ctx)
}
