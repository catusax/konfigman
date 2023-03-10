package tui

import (
	"errors"
	"fmt"

	"github.com/catusax/konfigman/pkg/kubeconfig"
	"github.com/manifoldco/promptui"
)

func SelectConfig() (string, error) {
	configs, err := kubeconfig.ListConfigs()
	if err != nil {
		return "", fmt.Errorf("get configs: %v", err)
	}

	if len(configs) == 0 {
		return "", errors.New("no config in registry")
	}

	prompt := promptui.Select{
		Label:    "Select a configuration",
		Items:    configs,
		HideHelp: true,
	}

	_, value, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("selection error: %v", err)
	}
	return value, nil
}
