package kubeconfig

import (
	"errors"
	"os"
	"path/filepath"
)

func GetCurrentConfig() (string, error) {
	target, err := os.Readlink(kubeconfig())
	if err != nil {
		return "", err
	}

	if filepath.Dir(target) == registryDir() {
		return filepath.Base(target), nil
	}
	return "", errors.New("not using registry")
}
