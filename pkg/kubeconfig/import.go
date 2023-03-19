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
	Path  string
	SSH   SSHConfig
	Host  string
	Name  string
	Force bool
}

type SSHConfig struct {
	URL          string
	IdentityFile string
	ConfigFile   string
	JumpHost     string
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

func ImportOptionWithSSH(ssh SSHConfig) ImportOption {
	return func(i *importOptions) {
		i.SSH = ssh
	}
}

func ImportOptionWithPath(path string) ImportOption {
	return func(i *importOptions) {
		i.Path = path
	}
}

func ImportOptionWithForce(force bool) ImportOption {
	return func(i *importOptions) {
		i.Force = force
	}
}

func ImportConfig(options ...ImportOption) (err error) {
	var option importOptions
	for _, o := range options {
		o(&option)
	}

	var config *clientcmdapi.Config
	if option.SSH.URL == "" {
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

		if option.SSH.URL != "" {
			url, _ := url.Parse(option.SSH.URL)

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
	if !os.IsNotExist(err) && !option.Force {
		return fmt.Errorf("%w,  %s", ErrConfAlreadyExists, option.Name)
	}

	err = clientcmd.WriteToFile(*config, filepath.Join(registryDir(), option.Name))
	if err != nil {
		return fmt.Errorf("failed to write configuration: %v", err)
	}

	return nil
}

func loadSSHConfig(ssh SSHConfig) (*clientcmdapi.Config, error) {
	params := append(buildSSHParams(ssh), ssh.URL, "kubectl config view --raw")
	command := exec.Command("ssh", params...)
	out, err := command.Output()
	if err != nil {
		return nil, fmt.Errorf("read remote config from ssh: %v", err)
	}

	return clientcmd.Load(out)
}

func buildSSHParams(ssh SSHConfig) (res []string) {
	if ssh.ConfigFile != "" {
		res = append(res, "-F", ssh.ConfigFile)
	}
	if ssh.IdentityFile != "" {
		res = append(res, "-i", ssh.IdentityFile)
	}
	if ssh.JumpHost != "" {
		res = append(res, "-J", ssh.JumpHost)
	}

	return
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
