package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

func UseConfig(name string) error {

	err := os.Remove(kubeconfig())
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	err = os.Symlink(filepath.Join(registryDir(), name), kubeconfig())
	if err != nil {
		return fmt.Errorf("link config %s %w", name, err)
	}
	return nil
}
