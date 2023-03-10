package kubeconfig

import "os"

func ListConfigs() ([]string, error) {
	entries, err := os.ReadDir(registryDir())
	if err != nil {
		return nil, err
	}

	var names []string
	for _, entry := range entries {
		if !entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	return names, nil
}
