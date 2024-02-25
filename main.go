package main

import (
	"fmt"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Clusters   []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			Server                   string `yaml:"server"`
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Users []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
			Token                 string `yaml:"token"`
		} `yaml:"user"`
	} `yaml:"users"`
	Contexts []struct {
		Name    string `yaml:"name"`
		Context struct {
			Cluster   string `yaml:"cluster"`
			User      string `yaml:"user"`
			Namespace string `yaml:"namespace"`
		} `yaml:"context"`
	} `yaml:"contexts"`
	CurrentContext string `yaml:"current-context"`
}

func main() {

	// read local kubeconfig and use the KubeConfig struct above
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	//fmt.Println(kubeconfig)
	// to unmarshal the kubeconfig file into the struct
	file, err := os.Open(kubeconfig)
	if err != nil {
		fmt.Println("Error opening kubeconfig file:", err)
		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	var config KubeConfig
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error decoding kubeconfig file:", err)
		return
	}

	for i, _ := range config.Contexts {
		cmd := exec.Command("kubectl", "get", "pods", "--context", config.Contexts[i].Name)
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Removing dead context: %s\n ", config.Contexts[i].Name)
			exec.Command("kubectl", "config", "delete-context", config.Contexts[i].Name).Run()
			exec.Command("kubectl", "config", "delete-cluster", config.Contexts[i].Context.Cluster).Run()
			exec.Command("kubectl", "config", "delete-user", config.Contexts[i].Context.User).Run()
		} else {
			fmt.Printf("Context is alive: %s\n ", config.Contexts[i].Name)
			//exec.Command("kubectl", "config", "use-context", config.Contexts[i].Name).Run()
		}

	}
}
