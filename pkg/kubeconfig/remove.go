package kubeconfig

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrConfigInUse = errors.New("config in use")

func RemoveConfig(name string, force bool) error {
	if !force {
		current, _ := GetCurrentConfig()
		if current == name {
			return ErrConfigInUse
		}
	}

	return os.Remove(filepath.Join(registryDir(), name))
}
