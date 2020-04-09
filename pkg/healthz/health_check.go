package healthz

import (
	"net/http"

	"k8s.io/klog"
)

func HandleHealthCheckRequest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte("{}")); err != nil {
		klog.Errorf("error writing health check response: %s", err)
	}
}
