package cloudprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/xiaods/k8e/pkg/version"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

var (
	InternalIPKey = version.Program + ".io/internal-ip"
	ExternalIPKey = version.Program + ".io/external-ip"
	HostnameKey   = version.Program + ".io/hostname"
)

func (k *k8e) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return cloudprovider.NotImplemented
}

func (k *k8e) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (k *k8e) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	return true, nil
}

func (k *k8e) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	if k.nodeInformerHasSynced == nil || !k.nodeInformerHasSynced() {
		return "", errors.New("Node informer has not synced yet")
	}

	_, err := k.nodeInformer.Lister().Get(string(nodeName))
	if err != nil {
		return "", fmt.Errorf("Failed to find node %s: %v", nodeName, err)
	}
	return string(nodeName), nil
}

func (k *k8e) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	return true, cloudprovider.NotImplemented
}

func (k *k8e) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	_, err := k.InstanceID(ctx, name)
	if err != nil {
		return "", err
	}
	return version.Program, nil
}

func (k *k8e) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	return "", cloudprovider.NotImplemented
}

func (k *k8e) NodeAddresses(ctx context.Context, name types.NodeName) ([]corev1.NodeAddress, error) {
	addresses := []corev1.NodeAddress{}
	if k.nodeInformerHasSynced == nil || !k.nodeInformerHasSynced() {
		return nil, errors.New("Node informer has not synced yet")
	}

	node, err := k.nodeInformer.Lister().Get(string(name))
	if err != nil {
		return nil, fmt.Errorf("Failed to find node %s: %v", name, err)
	}
	// check internal address
	if address := node.Annotations[InternalIPKey]; address != "" {
		for _, v := range strings.Split(address, ",") {
			addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: v})
		}
	} else if address = node.Labels[InternalIPKey]; address != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeInternalIP, Address: address})
	} else {
		logrus.Infof("Couldn't find node internal ip annotation or label on node %s", name)
	}

	// check external address
	if address := node.Annotations[ExternalIPKey]; address != "" {
		for _, v := range strings.Split(address, ",") {
			addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalIP, Address: v})
		}
	} else if address = node.Labels[ExternalIPKey]; address != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeExternalIP, Address: address})
	}

	// check hostname
	if address := node.Annotations[HostnameKey]; address != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeHostName, Address: address})
	} else if address = node.Labels[HostnameKey]; address != "" {
		addresses = append(addresses, corev1.NodeAddress{Type: corev1.NodeHostName, Address: address})
	} else {
		logrus.Infof("Couldn't find node hostname annotation or label on node %s", name)
	}

	return addresses, nil
}

func (k *k8e) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]corev1.NodeAddress, error) {
	return nil, cloudprovider.NotImplemented
}
