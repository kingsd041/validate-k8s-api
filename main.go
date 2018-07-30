package main

import (
	"fmt"

	"github.com/rancher/longhorn-manager/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	longhorn "github.com/rancher/longhorn-manager/k8s/pkg/apis/longhorn/v1alpha1"
	lhclientset "github.com/rancher/longhorn-manager/k8s/pkg/client/clientset/versioned"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	lhClient, err := lhclientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < 10000; i++ {
		settings := &longhorn.Setting{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(types.SettingNameStorageMinimalAvailablePercentage),
			},
			Setting: types.Setting{
				Value: "10",
			},
		}
		// create settings
		cs, err := lhClient.LonghornV1alpha1().Settings("longhorn-system").Create(settings)
		if err != nil {
			fmt.Printf("create settings error : %v \n", err)
		}
		fmt.Printf("create settings complete, result is: %+v \n", cs)
		s, err := lhClient.LonghornV1alpha1().Settings("longhorn-system").Get(string(types.SettingNameStorageMinimalAvailablePercentage), metav1.GetOptions{})
		if err != nil {
			fmt.Printf("get settings error : %v \n", err)
		}
		fmt.Printf("get settings after create, result is %+v \n", s)

		// update settings
		s.Value = "20"
		us, err := lhClient.LonghornV1alpha1().Settings("longhorn-system").Update(s)
		if err != nil {
			fmt.Printf("update settings error : %v \n", err)
		}
		fmt.Printf("update settings complete, result is %+v \n", us)
		aus, err := lhClient.LonghornV1alpha1().Settings("longhorn-system").Get(string(types.SettingNameStorageMinimalAvailablePercentage), metav1.GetOptions{})
		if err != nil {
			fmt.Printf("get settings error : %v \n", err)
		}
		fmt.Printf("get settings after update, result is %+v \n", aus)

		// delete settings
		err = lhClient.LonghornV1alpha1().Settings("longhorn-system").Delete(string(types.SettingNameStorageMinimalAvailablePercentage), &metav1.DeleteOptions{})
		if err != nil {
			fmt.Printf("delete settings error : %v \n", err)
		}
		sl, err := lhClient.LonghornV1alpha1().Settings("longhorn-system").List(metav1.ListOptions{})
		if err != nil {
			fmt.Printf("list settings error : %v \n", err)
		}
		fmt.Printf("list after delete, result is %+v \n", sl)
	}
}
