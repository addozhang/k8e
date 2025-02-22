package executor

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xiaods/k8e/pkg/cli/cmds"
	daemonconfig "github.com/xiaods/k8e/pkg/daemons/config"
	yaml2 "gopkg.in/yaml.v2"
	"k8s.io/apiserver/pkg/authentication/authenticator"
	"sigs.k8s.io/yaml"
)

var (
	executor Executor
)

type Executor interface {
	Bootstrap(ctx context.Context, nodeConfig *daemonconfig.Node, cfg cmds.Agent) error
	Kubelet(args []string) error
	KubeProxy(args []string) error
	APIServerHandlers() (authenticator.Request, http.Handler, error)
	APIServer(ctx context.Context, etcdReady <-chan struct{}, args []string) error
	Scheduler(apiReady <-chan struct{}, args []string) error
	ControllerManager(apiReady <-chan struct{}, args []string) error
	CurrentETCDOptions() (InitialOptions, error)
	ETCD(args ETCDConfig, extraArgs []string) error
	CloudControllerManager(ccmRBACReady <-chan struct{}, args []string) error
}

type ETCDConfig struct {
	InitialOptions      `json:",inline"`
	Name                string      `json:"name,omitempty"`
	ListenClientURLs    string      `json:"listen-client-urls,omitempty"`
	ListenMetricsURLs   string      `json:"listen-metrics-urls,omitempty"`
	ListenPeerURLs      string      `json:"listen-peer-urls,omitempty"`
	AdvertiseClientURLs string      `json:"advertise-client-urls,omitempty"`
	DataDir             string      `json:"data-dir,omitempty"`
	SnapshotCount       int         `json:"snapshot-count,omitempty"`
	ServerTrust         ServerTrust `json:"client-transport-security"`
	PeerTrust           PeerTrust   `json:"peer-transport-security"`
	ForceNewCluster     bool        `json:"force-new-cluster,omitempty"`
	HeartbeatInterval   int         `json:"heartbeat-interval"`
	ElectionTimeout     int         `json:"election-timeout"`
	Logger              string      `json:"logger"`
	LogOutputs          []string    `json:"log-outputs"`
}

type ServerTrust struct {
	CertFile       string `json:"cert-file"`
	KeyFile        string `json:"key-file"`
	ClientCertAuth bool   `json:"client-cert-auth"`
	TrustedCAFile  string `json:"trusted-ca-file"`
}

type PeerTrust struct {
	CertFile       string `json:"cert-file"`
	KeyFile        string `json:"key-file"`
	ClientCertAuth bool   `json:"client-cert-auth"`
	TrustedCAFile  string `json:"trusted-ca-file"`
}

type InitialOptions struct {
	AdvertisePeerURL string `json:"initial-advertise-peer-urls,omitempty"`
	Cluster          string `json:"initial-cluster,omitempty"`
	State            string `json:"initial-cluster-state,omitempty"`
}

func (e ETCDConfig) ToConfigFile(extraArgs []string) (string, error) {
	confFile := filepath.Join(e.DataDir, "config")
	bytes, err := yaml.Marshal(&e)
	if err != nil {
		return "", err
	}

	if len(extraArgs) > 0 {
		var s map[string]interface{}
		if err := yaml2.Unmarshal(bytes, &s); err != nil {
			return "", err
		}

		for _, v := range extraArgs {
			extraArg := strings.SplitN(v, "=", 2)
			// Depending on the argV, we have different types to handle.
			// Source: https://github.com/etcd-io/etcd/blob/44b8ae145b505811775f5af915dd19198d556d55/server/config/config.go#L36-L190 and https://etcd.io/docs/v3.5/op-guide/configuration/#configuration-file
			if len(extraArg) == 2 {
				key := strings.TrimLeft(extraArg[0], "-")
				lowerKey := strings.ToLower(key)
				var stringArr []string
				if i, err := strconv.Atoi(extraArg[1]); err == nil {
					s[key] = i
				} else if time, err := time.ParseDuration(extraArg[1]); err == nil && (strings.Contains(lowerKey, "time") || strings.Contains(lowerKey, "duration") || strings.Contains(lowerKey, "interval") || strings.Contains(lowerKey, "retention")) {
					// auto-compaction-retention is either a time.Duration or int, depending on version. If it is an int, it will be caught above.
					s[key] = time
				} else if err := yaml.Unmarshal([]byte(extraArg[1]), &stringArr); err == nil {
					s[key] = stringArr
				} else {
					switch strings.ToLower(extraArg[1]) {
					case "true":
						s[key] = true
					case "false":
						s[key] = false
					default:
						s[key] = extraArg[1]
					}
				}
			}
		}

		bytes, err = yaml2.Marshal(&s)
		if err != nil {
			return "", err
		}
	}

	if err := os.MkdirAll(e.DataDir, 0700); err != nil {
		return "", err
	}
	return confFile, ioutil.WriteFile(confFile, bytes, 0600)
}

func Set(driver Executor) {
	executor = driver
}

func Bootstrap(ctx context.Context, nodeConfig *daemonconfig.Node, cfg cmds.Agent) error {
	return executor.Bootstrap(ctx, nodeConfig, cfg)
}

func Kubelet(args []string) error {
	return executor.Kubelet(args)
}

func KubeProxy(args []string) error {
	return executor.KubeProxy(args)
}

func APIServerHandlers() (authenticator.Request, http.Handler, error) {
	return executor.APIServerHandlers()
}

func APIServer(ctx context.Context, etcdReady <-chan struct{}, args []string) error {
	return executor.APIServer(ctx, etcdReady, args)
}

func Scheduler(apiReady <-chan struct{}, args []string) error {
	return executor.Scheduler(apiReady, args)
}

func ControllerManager(apiReady <-chan struct{}, args []string) error {
	return executor.ControllerManager(apiReady, args)
}

func CurrentETCDOptions() (InitialOptions, error) {
	return executor.CurrentETCDOptions()
}

func ETCD(args ETCDConfig, extraArgs []string) error {
	return executor.ETCD(args, extraArgs)
}

func CloudControllerManager(ccmRBACReady <-chan struct{}, args []string) error {
	return executor.CloudControllerManager(ccmRBACReady, args)
}
