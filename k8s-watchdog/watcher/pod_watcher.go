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
)

func WatchPods(clientset *kubernetes.Clientset, namespaces []string) {
	if clientset == nil {
		log.Fatal("Clientset provided to WatchPods cannot be nil.")
	}
	log.Println("Starting Pod watchers for specified namespaces...")

	for _, ns := range namespaces {
		go func(namespace string) {
			log.Printf("Setting up watcher for namespace: %s\n", namespace)
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
				UpdateFunc: func(oldObj, newObj interface{}) {
					oldPod := oldObj.(*corev1.Pod)
					newPod := newObj.(*corev1.Pod)
					if oldPod.ResourceVersion != newPod.ResourceVersion {
						fmt.Printf("[~] Pod updated in %s: %s\n", namespace, newPod.GetName())
					}
				},
			})

			stopCh := make(chan struct{})
			defer close(stopCh)

			go factory.Start(stopCh)
			factory.WaitForCacheSync(stopCh)
			log.Printf("Cache synced for namespace: %s. Ready to watch events.\n", namespace)

			<-stopCh
		}(ns)
	}

	select {}
}
