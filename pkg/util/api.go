package util

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rancher/wrangler/pkg/merr"
	"github.com/rancher/wrangler/pkg/schemes"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	coregetter "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
)

// This sets a default duration to wait for the apiserver to become ready. This is primarily used to
// block startup of agent supervisor controllers until the apiserver is ready to serve requests, in the
// same way that the apiReady channel is used in the server packages, so it can be fairly long. It must
// be at least long enough for downstream projects like RKE2 to start the apiserver in the background.
const DefaultAPIServerReadyTimeout = 15 * time.Minute

func GetAddresses(endpoint *v1.Endpoints) []string {
	serverAddresses := []string{}
	if endpoint == nil {
		return serverAddresses
	}
	for _, subset := range endpoint.Subsets {
		var port string
		if len(subset.Ports) > 0 {
			port = strconv.Itoa(int(subset.Ports[0].Port))
		}
		if port == "" {
			port = "443"
		}
		for _, address := range subset.Addresses {
			serverAddresses = append(serverAddresses, net.JoinHostPort(address.IP, port))
		}
	}
	return serverAddresses
}

// WaitForAPIServerReady waits for the API Server's /readyz endpoint to report "ok" with timeout.
// This is modified from WaitForAPIServer from the Kubernetes controller-manager app, but checks the
// readyz endpoint instead of the deprecated healthz endpoint, and supports context.
func WaitForAPIServerReady(ctx context.Context, client clientset.Interface, timeout time.Duration) error {
	var lastErr error
	restClient := client.Discovery().RESTClient()

	err := wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		healthStatus := 0
		result := restClient.Get().AbsPath("/readyz").Do(ctx).StatusCode(&healthStatus)
		if rerr := result.Error(); rerr != nil {
			if errors.Is(rerr, context.Canceled) {
				return false, rerr
			}
			lastErr = errors.Wrap(rerr, "failed to get apiserver /readyz status")
			return false, nil
		}
		if healthStatus != http.StatusOK {
			content, _ := result.Raw()
			lastErr = fmt.Errorf("APIServer isn't ready: %v", string(content))
			logrus.Warnf("APIServer isn't ready yet: %v. Waiting a little while.", string(content))
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		return merr.NewErrors(err, lastErr)
	}

	return nil
}

func BuildControllerEventRecorder(k8s clientset.Interface, controllerName string) record.EventRecorder {
	logrus.Infof("Creating %s event broadcaster", controllerName)
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(&coregetter.EventSinkImpl{Interface: k8s.CoreV1().Events(metav1.NamespaceSystem)})
	nodeName := os.Getenv("NODE_NAME")
	return eventBroadcaster.NewRecorder(schemes.All, v1.EventSource{Component: controllerName, Host: nodeName})
}
