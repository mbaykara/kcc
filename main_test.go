package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadKubeConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "kubeconfig_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir) // Clean up after the test

	tmpKubeConfig := filepath.Join(tmpDir, ".kube", "config")

	// Ensure the .kube directory exists
	if err := os.MkdirAll(filepath.Dir(tmpKubeConfig), 0755); err != nil {
		t.Fatalf("Failed to create .kube directory: %v", err)
	}

	// Sample kubeconfig data
	kubeConfigContent := `
apiVersion: v1
kind: Config
clusters:
- name: test-cluster
  cluster:
    server: "https://example.com"
    certificate-authority-data: "test"
users:
- name: test-user
  user:
    client-certificate-data: "test"
    client-key-data: "test"
contexts:
- name: test-context
  context:
    cluster: "test-cluster"
    user: "test-user"
current-context: "test-context"
`
	// Write the sample data to the temporary kubeconfig file
	if err := os.WriteFile(tmpKubeConfig, []byte(kubeConfigContent), 0644); err != nil {
		t.Fatalf("Failed to write kubeconfig file: %v", err)
	}

	// Override HOME environment variable to point to the temporary directory
	originalHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", tmpDir); err != nil {
		t.Fatalf("Failed to set HOME environment variable: %v", err)
	}
	defer os.Setenv("HOME", originalHome)

	expected := KubeConfig{
		APIVersion: "v1",
		Kind:       "Config",
		Clusters: []struct {
			Name    string  `yaml:"name"`
			Cluster Cluster `yaml:"cluster"`
		}{
			{
				Name: "test-cluster",
				Cluster: Cluster{
					Server:                   "https://example.com",
					CertificateAuthorityData: "test",
				},
			},
		},
		Users: []struct {
			Name string `yaml:"name"`
			User User   `yaml:"user"`
		}{
			{
				Name: "test-user",
				User: User{
					ClientCertificateData: "test",
					ClientKeyData:         "test",
				},
			},
		},
		Contexts: []struct {
			Name    string  `yaml:"name"`
			Context Context `yaml:"context"`
		}{
			{
				Name: "test-context",
				Context: Context{
					Cluster:   "test-cluster",
					User:      "test-user",
					Namespace: "",
				},
			},
		},
		CurrentContext: "test-context",
	}

	result, err := readKubeConfig()
	if err != nil {
		t.Fatalf("readKubeConfig() returned an error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("readKubeConfig() = %v, want %v", result, expected)
	}
}
