package activator

import (
	"context"
	"sync"
	"time"

	k8s "github.com/deislabs/osiris/pkg/kubernetes"
	"k8s.io/klog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type deploymentActivation struct {
	readyAppPodIPs map[string]struct{}
	endpoints      *corev1.Endpoints
	lock           sync.Mutex
	successCh      chan string
	timeoutCh      chan struct{}
}

func (d *deploymentActivation) watchForCompletion(
	kubeClient kubernetes.Interface,
	app *app,
	appPodSelector labels.Selector,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Watch the pods managed by this deployment
	podsInformer := k8s.PodsIndexInformer(
		kubeClient,
		app.namespace,
		nil,
		appPodSelector,
	)
	podsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: d.syncPod,
		UpdateFunc: func(_, newObj interface{}) {
			d.syncPod(newObj)
		},
		DeleteFunc: d.syncPod,
	})
	// Watch the corresponding endpoints resource for this service
	endpointsInformer := k8s.EndpointsIndexInformer(
		kubeClient,
		app.namespace,
		fields.OneTermEqualSelector(
			"metadata.name",
			app.serviceName,
		),
		nil,
	)
	endpointsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: d.syncEndpoints,
		UpdateFunc: func(_, newObj interface{}) {
			d.syncEndpoints(newObj)
		},
	})
	go podsInformer.Run(ctx.Done())
	go endpointsInformer.Run(ctx.Done())
	timer := time.NewTimer(2 * time.Minute)
	defer timer.Stop()
	for {
		select {
		case <-d.successCh:
			return
		case <-timer.C:
			klog.Errorf(
				"Activation of deployment %s in namespace %s timed out",
				app.deploymentName,
				app.namespace,
			)
			close(d.timeoutCh)
			return
		}
	}
}

func (d *deploymentActivation) syncPod(obj interface{}) {
	d.lock.Lock()
	defer d.lock.Unlock()
	pod := obj.(*corev1.Pod)
	klog.Infof("so we got to syncPod with this %v:", pod.Status.Conditions)
	var ready bool
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			if condition.Status == corev1.ConditionTrue {
				klog.Infof("Pod is Ready")
				ready = true
			}
			break
		}
	}
	klog.Infof("we are %s, getting further and we need to return this IP %v", ready, pod.Status.PodIP)
	// Keep track of which pods are ready
	if ready {
		d.readyAppPodIPs[pod.Status.PodIP] = struct{}{}
	} else {
		delete(d.readyAppPodIPs, pod.Status.PodIP)
	}
	d.checkActivationComplete()
}

func (d *deploymentActivation) syncEndpoints(obj interface{}) {
	klog.Infof("so we got to syncEndpoints with this %v:", obj)
	d.lock.Lock()
	defer d.lock.Unlock()
	d.endpoints = obj.(*corev1.Endpoints)
	d.checkActivationComplete()
}

func (d *deploymentActivation) checkActivationComplete() {
	klog.Infof("And we reached checkActivationComplete() with this endpoints: %v", d.endpoints)
	if d.endpoints != nil {
		for _, subset := range d.endpoints.Subsets {
			for _, address := range subset.Addresses {
				if _, ok := d.readyAppPodIPs[address.IP]; ok {
					klog.Infof("App pod with ip %s is in service", address.IP)
					close(d.successCh)
					return
				}
			}
		}
	}
}
