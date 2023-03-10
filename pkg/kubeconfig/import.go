package kubeconfig

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type importOptions struct {
	Path string
	SSH  string
	Host string
	Name string
}

type ImportOption func(i *importOptions)

var ErrConfAlreadyExists = errors.New("config already exists")

func ImportOptionWithHost(host string) ImportOption {
	return func(i *importOptions) {
		i.Host = host
	}
}

func ImportOptionWithName(name string) ImportOption {
	return func(i *importOptions) {
		i.Name = name
	}
}

func ImportOptionWithSSH(ssh string) ImportOption {
	return func(i *importOptions) {
		i.SSH = ssh
	}
}

func ImportOptionWithPath(path string) ImportOption {
	return func(i *importOptions) {
		i.Path = path
	}
}

func ImportConfig(options ...ImportOption) (err error) {
	var option importOptions
	for _, o := range options {
		o(&option)
	}

	var config *clientcmdapi.Config
	if option.SSH == "" {
		config, err = clientcmd.LoadFromFile(option.Path)
		if err != nil {
			return err
		}
	} else {
		config, err = loadSSHConfig(option.SSH)
		if err != nil {
			return err
		}
	}

	for i := range config.Clusters {

		if option.SSH != "" {
			url, _ := url.Parse(option.SSH)

			serverurl, err := url.Parse(config.Clusters[i].Server)
			if err != nil {
				return fmt.Errorf("failed to parse server url %s: %v", config.Clusters[i].Server, err)
			}

			if serverurl.Hostname() == "127.0.0.1" {
				err = replaceHostName(url.Hostname(), config.Clusters[i])
				if err != nil {
					return fmt.Errorf("replace hostname %s:%v", url.Hostname(), err)
				}
			}
		}

		if option.Host != "" {
			err = replaceHostName(option.Host, config.Clusters[i])
			if err != nil {
				return fmt.Errorf("replace hostname %s:%v", option.Host, err)
			}
		}
	}

	_, err = os.Stat(filepath.Join(registryDir(), option.Name))
	if !os.IsNotExist(err) {
		return fmt.Errorf("%w,  %s", ErrConfAlreadyExists, option.Name)
	}

	err = clientcmd.WriteToFile(*config, filepath.Join(registryDir(), option.Name))
	if err != nil {
		return fmt.Errorf("failed to write configuration: %v", err)
	}

	return nil
}

func loadSSHConfig(ssh string) (*clientcmdapi.Config, error) {
	command := exec.Command("ssh", ssh, "kubectl config view --raw")
	out, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("read remote config from ssh: %v", err)
	}

	return clientcmd.Load(out)
}

func replaceHostName(hostName string, cluster *clientcmdapi.Cluster) error {
	serverurl, err := url.Parse(cluster.Server)
	if err != nil {
		return fmt.Errorf("failed to parse server url %s: %v", cluster.Server, err)
	}
	oldHost := serverurl.Hostname()

	cluster.TLSServerName = serverurl.Hostname()

	serverurl.Host = hostName + ":" + serverurl.Port()

	cluster.Server = serverurl.String()
	cluster.TLSServerName = oldHost
	return nil
}
