package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
)

type UseOptions struct {
	Context   string
	Namespace string
	Cluster   string
	User      string
}

func UseConfig(name string, options UseOptions) error {
	sourceFile := filepath.Join(registryDir(), name)

	err := os.Remove(kubeconfig())
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}

	config, err := clientcmd.LoadFromFile(sourceFile)
	if err != nil {
		return fmt.Errorf("read config %s: %w", name, err)
	}

	currentContext := config.CurrentContext
	if options.Context != "" {
		currentContext = options.Context
	}
	if options.Namespace != "" {
		config.Contexts[currentContext].Namespace = options.Namespace
	}
	if options.User != "" {
		config.Contexts[currentContext].AuthInfo = options.User
	}

	err = clientcmd.WriteToFile(*config, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to write configuration: %v", err)
	}

	err = os.Symlink(filepath.Join(registryDir(), name), kubeconfig())
	if err != nil {
		return fmt.Errorf("link config %s: %w", name, err)
	}

	return nil
}
