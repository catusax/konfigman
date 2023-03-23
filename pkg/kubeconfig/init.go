package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"
)

func InitRegistry() error {
	_, err := os.ReadDir(registryDir())
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(registryDir(), 0755)
			if err != nil {
				return fmt.Errorf("create registry dir failed: %v", err)
			}

			_, err := os.Stat(kubeconfig())
			if err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("read current kubeconfig failed ~/.kube/config : %v", err)
				}
			} else {
				err = os.Rename(kubeconfig(), filepath.Join(registryDir(), "default"))
				if err != nil {
					return fmt.Errorf("import kubeconfig failed: %v", err)
				}
			}
		}
	}

	return nil
}
