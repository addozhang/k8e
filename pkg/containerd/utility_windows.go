// +build windows

package containerd

import (
	"github.com/pkg/errors"
	util2 "github.com/xiaods/k8e/pkg/util"
)

func OverlaySupported(root string) error {
	return errors.Wrapf(util2.ErrUnsupportedPlatform, "overlayfs is not supported")
}

func FuseoverlayfsSupported(root string) error {
	return errors.Wrapf(util2.ErrUnsupportedPlatform, "fuse-overlayfs is not supported")
}