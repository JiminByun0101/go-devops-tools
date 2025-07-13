package watcher

import (
	"fmt"
	"log"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func WatchPods(namespaces []string) {
	kubeconfig := clientcmd.RecommendedHomeFile

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	for _, ns := range namespaces {
		go func(namespace string) {
			factory := informers.NewSharedInformerFactoryWithOptions(
				clientset,
				time.Minute,
				informers.WithNamespace(namespace),
				informers.WithTweakListOptions(func(opt *metav1.ListOptions) {
					opt.FieldSelector = fields.Everything().String()
				}),
			)

			informer := factory.Core().V1().Pods().Informer()

			informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					pod := obj.(*corev1.Pod)
					fmt.Printf("[+] Pod added: %s%s\n", namespace, pod.GetName())
				},
				DeleteFunc: func(obj interface{}) {
					pod := obj.(*corev1.Pod)
					fmt.Printf("[-] Pod deleted: %s%s\n", namespace, pod.GetName())
				},
			})

			stopCh := make(chan struct{})
			defer close(stopCh)

			go factory.Start(stopCh)
			factory.WaitForCacheSync(stopCh)

			<-stopCh
		}(ns)
	}

	select {}
}
