package kubeconfig

import (
	"os"
	"path/filepath"
)

func kubeconfig() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}

func registryDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "konfigman-registry")
}
