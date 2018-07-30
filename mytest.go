package main

import (
	"fmt"
	"time"
	"context"
	"flag"
	"strings"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/transport"
	"github.com/rancher/longhorn-manager/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	longhorn "github.com/rancher/longhorn-manager/k8s/pkg/apis/longhorn/v1alpha1"
	lhclientset "github.com/rancher/longhorn-manager/k8s/pkg/client/clientset/versioned"
)

var (
	testk8s bool
	testetcd bool
	certFile string
	keyFile string
	endpoints string
)

func main() {
	flag.Parse()

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	if testk8s {
		testK8s(config)
	} else if testetcd {
		testEtcd(config)
	}

}

func init() {
	flag.StringVar(&certFile, "certFile", "", "cert file for etcd server")
	flag.StringVar(&keyFile, "keyFile", "", "key file for etcd server")
	flag.StringVar(&endpoints, "endpoints", "", "etcd endpoints, separate with ,")
	flag.BoolVar(&testk8s, "testk8s", true, "enable test k8s api")
	flag.BoolVar(&testetcd, "testetcd", false, "enable test etcd api")
}

func testK8s(config *rest.Config) {
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

func testEtcd(config *rest.Config) {
	tlsInfo := transport.TLSInfo{
		CertFile: certFile,
		KeyFile:  keyFile,
		TrustedCAFile: config.CAFile,
	}
	tlsConfig, err := tlsInfo.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	etcdEndpoints := strings.Split(endpoints, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
		TLS: tlsConfig,
	})

	for i := 0; i < 1000; i++ {
		// create new key, value
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		rs, err := cli.Put(ctx, "/longhorn/test/", "test_Save_value")
		if err != nil {
			fmt.Printf("create through etcd error : %v \n", err)
		}
		fmt.Printf("create complete, result is: %+v \n", rs.PrevKv)
		cancel()
		// get after create
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		resp, err := cli.Get(ctx, "/longhorn/test/")
		if err != nil {
			fmt.Printf("get after create through etcd error : %v \n", err)
		}
		fmt.Printf("get after create, result is :%+v \n", resp.Kvs)
		cancel()
		// update key, value
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		us, err := cli.Put(ctx, "/longhorn/test/", "new_test_value")
		if err != nil {
			fmt.Printf("update through etcd error : %v \n", err)
		}
		fmt.Printf("update complete, result is: %+v \n", us.PrevKv)
		cancel()
		// get after update
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		ngs, err := cli.Get(ctx, "/longhorn/test")
		if err != nil {
			fmt.Printf("get after update through etcd error : %v \n", err)
		}
		fmt.Printf("get after create, result is :%+v \n", ngs.Kvs)
		cancel()
		//delete key, value
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		_, err = cli.Delete(ctx, "/longhorn/test/")
		if err != nil {
			fmt.Printf("delete through etcd error : %v \n", err)
		}
		cancel()
		// list after delete
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		gr, err := cli.Get(ctx, "/")
		if err != nil {
			fmt.Printf("get after delete through etcd error : %v \n", err)
		}
		fmt.Printf("get after create, result is :%+v \n", gr.Kvs)
		cancel()
	}
}
