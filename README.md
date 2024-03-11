# KCC
KCC (KubeConfigCleaner) is a lightweight tool designed to clean up your Kubernetes configuration files (kubeconfig) by removing dead or obsolete contexts. It helps keep your ~/.kube/config file organized and clutter-free.

Features
* Context Cleanup: Identify and remove inactive or obsolete contexts from your Kubernetes configuration file.
* User-Friendly Interface: Simple command-line interface for ease of use.


## Build locally

```bash
git clone https://github.com/mbaykara/kcc.git
cd kcc && go build -o kcc main.go
mv kcc /usr/local/bin
```

