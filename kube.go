package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	v2 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var email = " "
	var clientset *kubernetes.Clientset
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}
	clientset, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		config, err := rest.InClusterConfig()
		if err == nil {
			panic(err.Error())
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Printf("error getting kubernetes config: %v\n", err)
			os.Exit(1)
		}
	}

	// An empty string returns all namespaces
	//namespace := ""
	//pods, err := ListPods(namespace, clientset)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// // for _, pod := range pods.Items {
	// // 	fmt.Printf("Pod name: %v\n", pod.Name)
	// // }

	// for _, pod := range pods.Items {
	// 	if pod.Status.ContainerStatuses[0].State.Waiting != nil && pod.Status.ContainerStatuses[0].State.Waiting.Reason == "CrashLoopBackOff" && pod.Status.ContainerStatuses[0].RestartCount > 0 {
	// 		restartCount := pod.Status.ContainerStatuses[0].RestartCount
	// 		tempNamespace := pod.Namespace
	// 		podNamespace, err := getNamespaceLabels(tempNamespace, clientset)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			os.Exit(1)
	// 		}
	// 		for _, namespace := range podNamespace.Items {
	// 			if tempNamespace == namespace.Name {
	// 				for key, value := range namespace.Labels {
	// 					if key == "primary-owner" {
	// 						//fmt.Println(pod.Name, "in", pod.Namespace, "Namespace is in Crashloop Backoff")
	// 						email := strings.ReplaceAll(value, "_", "@")
	// 						fmt.Println(pod.Name, "in", pod.Namespace, "Namespace is in Crashloop Backoff with a restart count of ", restartCount, "and primary owner is", email)
	// 						//scaleDeployments(tempNamespace, clientset, "nginx")
	// 					}

	// 				}
	// 			}

	// 		}
	// 		// fmt.Println("Namespace Labels", podNamespace.Items[].Name)
	// 	}
	// }

	// var message string
	// if namespace == "" {
	// 	message = "Total Pods in all namespaces"
	// } else {
	// 	message = fmt.Sprintf("Total Pods in namespace `%s`", namespace)
	// }
	// fmt.Printf("%s %d\n", message, len(pods.Items))

	//ListNamespaces function call returns a list of namespaces in the kubernetes cluster

	namespaces, err := ListNamespaces(clientset)
	if err != nil {
		fmt.Println(err.Error)
		os.Exit(1)
	}
	for _, namespace := range namespaces.Items {
		if !strings.Contains(namespace.Name, "openshift") {
			//fmt.Fprintln(os.Stdout, []any{"deployments in namespace", namespace.Name}...)
			deployments, err := Listdeployments(namespace.Name, clientset)
			if err != nil {
				fmt.Fprintln(os.Stdout, []any{err.Error}...)
				os.Exit(1)
			}

			for _, deployments := range deployments.Items {
				//fmt.Fprintln(os.Stdout, []any{deployments.Name}...)
				tempnamespace := namespace.Name
				for key, value := range namespace.Labels {
					if key == "primary-owner" {
						//fmt.Println(pod.Name, "in", pod.Namespace, "Namespace is in Crashloop Backoff")
						email = strings.ReplaceAll(value, "_", "@")
						//scaleDeployments(tempNamespace, clientset, "nginx")
					}

				}
				tempdeployment := deployments.Name
				set := labels.Set(deployments.Spec.Selector.MatchLabels)
				pods, err := ListPods(namespace.Name, clientset, set.AsSelector().String())
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				for _, pod := range pods.Items {
					if pod.Status.ContainerStatuses[0].State.Waiting != nil && pod.Status.ContainerStatuses[0].State.Waiting.Reason == "CrashLoopBackOff" && pod.Status.ContainerStatuses[0].RestartCount > 0 {

						fmt.Println("namespace", tempnamespace, "deploymentname", tempdeployment, "podName", pod.Name, "email", email)
						//scaleDeployments(tempnamespace, clientset, tempdeployment)
					}
				}

			}
		}
	}
	//fmt.Printf("Total namespaces: %d\n", len(namespaces.Items))
}

func Listdeployments(namespace string, client kubernetes.Interface) (*v2.DeploymentList, error) {
	deployments, err := client.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting pods: %v\n", err)
		return nil, err
	}
	return deployments, nil
}

func ListPods(namespace string, client kubernetes.Interface, selectors string) (*v1.PodList, error) {
	listOptions := metav1.ListOptions{LabelSelector: selectors}
	pods, err := client.CoreV1().Pods(namespace).List(context.Background(), listOptions)
	if err != nil {
		err = fmt.Errorf("error getting pods: %v\n", err)
		return nil, err
	}
	return pods, nil
}

func ListNamespaces(client kubernetes.Interface) (*v1.NamespaceList, error) {
	fmt.Println("Get Kubernetes Namespaces")
	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting namespaces: %v\n", err)
		return nil, err
	}
	return namespaces, nil
}

func getNamespaceLabels(namespace string, client kubernetes.Interface) (*v1.NamespaceList, error) {
	namespaces, err := client.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		err = fmt.Errorf("error getting namespaces: %v\n", err)
		return nil, err
	}
	return namespaces, nil
}

func scaleDeployments(namespace string, client kubernetes.Interface, deploymentname string) error {
	s, err := client.AppsV1().
		Deployments(namespace).
		GetScale(context.TODO(), deploymentname, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}

	sc := *s
	sc.Spec.Replicas = 0

	us, err := client.AppsV1().
		Deployments(namespace).
		UpdateScale(context.TODO(),
			deploymentname, &sc, metav1.UpdateOptions{})
	_ = us
	if err != nil {
		log.Fatal(err)
	}
	return err
}
