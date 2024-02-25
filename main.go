package main

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type Cluster struct {
	Server                   string `yaml:"server"`
	CertificateAuthorityData string `yaml:"certificate-authority-data"`
}

type User struct {
	ClientCertificateData string `yaml:"client-certificate-data"`
	ClientKeyData         string `yaml:"client-key-data"`
	Token                 string `yaml:"token"`
}

type Contexts struct {
	Context Context `yaml:"context"`
	Name    string  `yaml:"name"`
}

type Context struct {
	Cluster   string `yaml:"cluster"`
	User      string `yaml:"user"`
	Namespace string `yaml:"namespace"`
}

type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Clusters   []struct {
		Name    string  `yaml:"name"`
		Cluster Cluster `yaml:"cluster"`
	} `yaml:"clusters"`
	Users []struct {
		Name string `yaml:"name"`
		User User   `yaml:"user"`
	} `yaml:"users"`
	Contexts []struct {
		Name    string  `yaml:"name"`
		Context Context `yaml:"context"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
}

func main() {
	config, err := readKubeConfig()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, context := range config.Contexts {
		if !checkContext(context.Name) {
			removeDeadContext(Contexts{
				Name:    context.Name,
				Context: context.Context,
			})
		} else {
			fmt.Printf("Context is alive: %s\n", context.Name)
		}
	}
}

func readKubeConfig() (KubeConfig, error) {
	var config KubeConfig
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	file, err := os.Open(kubeconfig)
	if err != nil {
		return config, fmt.Errorf("error opening kubeconfig file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("error decoding kubeconfig file: %w", err)
	}

	return config, nil
}

func checkContext(contextName string) bool {
	cmd := exec.Command("kubectl", "get", "pods", "--context", contextName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Removing dead context: %s\n", contextName)
		return false
	}
	return true
}

func removeDeadContext(context Contexts) {
	exec.Command("kubectl", "config", "delete-context", context.Name).Run()
	exec.Command("kubectl", "config", "delete-cluster", context.Context.Cluster).Run()
	exec.Command("kubectl", "config", "delete-user", context.Context.User).Run()
}
