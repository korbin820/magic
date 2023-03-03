package k8s

import (
	"context"
	"flag"
	"path/filepath"
	"strings"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var clientset *kubernetes.Clientset

func init() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func GetConfig(key string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cm, err := clientset.CoreV1().ConfigMaps("monitor").Get(ctx, "mtagsetting", v1.GetOptions{})
	defer cancel()
	if err != nil {
		panic("获取集群configmap异常!")
	}

	val, judeg := cm.Data[key]
	if !judeg {
		return ""
	}

	return strings.TrimSpace(strings.Replace(val, "\n", "", -1))
}
