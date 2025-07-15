package main

import (
	"fmt"
	"log"

	"github.com/JiminByun0101/go-devops-tools/k8s-watchdog/config"
	"github.com/JiminByun0101/go-devops-tools/k8s-watchdog/pkg/kube"
	"github.com/JiminByun0101/go-devops-tools/k8s-watchdog/watcher"
)

func main() {
	fmt.Println("K8s Watchdog starting...")

	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Faild to load config: %v", err)
	}

	fmt.Printf("Watching resources: %v in namespaces: %v\n", cfg.Watch.Resources, cfg.Watch.Namespaces)

	clientset, err := kube.GetClientSet()
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	watcher.WatchPods(clientset, cfg.Watch.Namespaces)

}
