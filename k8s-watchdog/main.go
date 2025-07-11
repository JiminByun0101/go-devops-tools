package main

import (
	"fmt"
	"log"

	"github.com/JiminByun0101/go-devops-tools/k8s-watchdog/config"
)

func main() {
	fmt.Println("K8s Watchdog starting...")

	cfg, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Faild to load config: %v", err)
	}

	fmt.Printf("Watching resources: %v in namespaces: %v\n", cfg.Watch.Resources, cfg.Watch.Namespaces)

}
