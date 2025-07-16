# k8s-watchdog

A simple Go-based Kubernetes Pod watcher application leveraging `client-go` informers. This project serves as a hands-on guide demonstrating foundational concepts for building Kubernetes tools, including configuration management and flexible deployment strategies.

**📖 For a more detailed, step-by-step guide on how this project was built, please refer to the accompanying blog post on Medium:**
[Getting Started with client-go: Building a Kubernetes Pod Watcher in Go](https://jiminbyun.medium.com/getting-started-with-client-go-building-a-kubernetes-pod-watcher-in-go-caa2be8623eb)

## ✨ Demonstrated Concepts

This repository illustrates key `client-go` and Kubernetes deployment concepts, including:

* **Real-time Kubernetes Event Monitoring:** How to set up a basic Pod watcher using `client-go` informers to observe creation and deletion events.
* **External Configuration Management:** Implementing flexible application configuration using `config.yaml` and the `viper` library.
* **Dynamic Kubernetes Client Setup:** Handling both out-of-cluster (via `kubeconfig`) and in-cluster (via `ServiceAccount` tokens) Kubernetes API connectivity.
* **Application Containerization:** Packaging a Go application into a Docker image for Kubernetes deployment.
* **Standard Kubernetes Deployment:** Utilizing common Kubernetes manifests (Deployment, ConfigMap, ServiceAccount, ClusterRole, ClusterRoleBinding) for in-cluster application setup and permissions.


## 📋 Prerequisites

Before you begin, ensure you have the following installed on your machine:

* **Go:** Version 1.18 or higher ([https://golang.org/doc/install](https://golang.org/doc/install))
* **Docker:** For building container images ([https://docs.docker.com/get-docker/](https://docs.docker.com/get-docker/))
* **kubectl:** Kubernetes command-line tool ([https://kubernetes.io/docs/tasks/tools/install-kubectl/](https://kubernetes.io/docs/tasks/tools/install-kubectl/))
* **Minikube (Recommended):** For a local Kubernetes cluster ([https://minikube.sigs.k8s.io/docs/start/](https://minikube.sigs.k8s.io/docs/start/)) or any other Kubernetes cluster you have access to.

## 📁 Project Structure

```
.
├── Dockerfile                  \# Defines the Docker image build process
├── go.mod                      \# Go module definition
├── go.sum                      \# Go module checksums
├── main.go                     \# Main application entry point
├── config.yaml                 \# Default configuration file for namespaces and resources
├── k8s-watcher-deployment.yaml \# Kubernetes manifests for in-cluster deployment (RBAC, ConfigMap, Deployment)
├── config/                     \# Configuration package
│   └── config.go               \# Handles loading configuration via Viper
├── pkg/                        \# Utility packages
│   └── kube/                   \# Kubernetes client utilities
│       └── client.go           \# Provides Kubernetes Clientset (in-cluster/out-of-cluster)
└── watcher/                    \# Watcher logic package
    └── pod\_watcher.go         \# Implements the Pod watching using client-go informers
````

## 🚀 Getting Started

Follow these steps to get your `k8s-watchdog` running.

### 1. Clone the Repository

```bash
git clone https://github.com/JiminByun0101/go-devops-tools.git
cd go-devops-tools/k8s-watchdog
````

### 2\. Install Go Modules

Navigate to the project root and download the necessary Go modules:

```bash
go mod tidy
```

### 3\. Configure the Watcher

The application reads its configuration from `config.yaml`. A basic `config.yaml` is provided at the root of the repository.

```yaml
# config.yaml
watch:
  namespaces: ["default", "kube-system"] # Namespaces to watch (e.g., "default", "all" for all)
  resources: ["pods"]                    # Resources to watch (e.g., "pods", "deployments")
```

You can modify `namespaces` and `resources` to suit your needs.

### 4\. Run Locally (Out-of-Cluster Test)

To test the watcher on your local machine, ensure you have `kubectl` configured to connect to your Kubernetes cluster (e.g., Minikube is running and `kubectl config use minikube`).

```bash
go run main.go
```

You should see output similar to this, indicating it's using your local kubeconfig:

```
K8s Watchdog starting...
Watching resources: [pods] in namespaces: [default kube-system]
2025/07/15 17:15:10 Using out-of-cluster Kubernetes configuration. (kubeconfig: /home/youruser/.kube/config).
2025/07/15 17:15:10 Starting Pod watchers for specified namespaces...
# ... watch events will follow ...
```

Try creating or deleting a Pod in your `default` namespace (e.g., `kubectl run nginx --image=nginx --restart=Never -- sleep 3600`) to see it in action.

### 5\. Build and Deploy to Kubernetes (In-Cluster Test)

This section demonstrates how to containerize your application and deploy it to a Kubernetes cluster, allowing it to use **in-cluster configuration**.

#### 5.1. Build the Docker Image

If using Minikube, configure your shell to use Minikube's Docker daemon. This command ensures Docker commands build directly into Minikube's Docker environment:

```bash
minikube start # If your Minikube cluster isn't running
eval $(minikube docker-env)
```

Now, build the Docker image for your application:

```bash
docker build -t k8s-watchdog:latest .
```

#### 5.2. Deploy to Kubernetes

The `k8s-watcher-deployment.yaml` file contains all the necessary Kubernetes manifests:

  * A `ServiceAccount` to provide an identity for your Pod.
  * A `ClusterRole` and `ClusterRoleBinding` to grant the `ServiceAccount` permissions to `list` and `watch` Pods across namespaces (essential for watching `kube-system`).
  * A `ConfigMap` to provide the `config.yaml` to your running Pod.
  * A `Deployment` to manage your application Pod.

Apply the Kubernetes manifests to your cluster:

```bash
kubectl apply -f k8s-watcher-deployment.yaml
```

#### 5.3. Verify In-Cluster Deployment and Logs

Wait for your Pod to be running. You can check its status with:

```bash
kubectl get pods -n default
# Look for a Pod named like k8s-watchdog-deployment-xxxx-yyyy
```

Once running, view the logs from your deployed Pod. This is how you'll observe its behavior when running inside the cluster:

```bash
kubectl logs -f deployment/k8s-watchdog-deployment -n default
```

You should see output similar to this, confirming it's using in-cluster configuration and watching events:

```
K8s Watchdog starting...
Watching resources: [pods] in namespaces: [default kube-system]
2025/07/15 17:18:26 Using in-cluster Kubernetes configuration.
# ... cache sync and watch events will follow for both default and kube-system ...
```

Try creating or deleting Pods in your `default` namespace (e.g., `kubectl run testpod --image=busybox --restart=Never -- sleep 3600`) and observe the logs from the `k8s-watchdog` Pod.

## 🧹 Cleanup

To remove the deployed application from your cluster:

```bash
kubectl delete -f k8s-watcher-deployment.yaml
```

If you started Minikube for this tutorial, you can stop it:

```bash
minikube stop
```