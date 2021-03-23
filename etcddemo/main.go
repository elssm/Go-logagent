package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main()  {
	cli,err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: time.Second*5,
	})
	if err!=nil {
		fmt.Printf("connect to etcd failed , err: %v",err)
		return
	}
	defer cli.Close()

	ctx,cancel := context.WithTimeout(context.Background(),time.Second)
	str := `[{"path":"/Users/caoyifan/Desktop/s4.log","topic":"s4_log"},{"path":"/Users/caoyifan/Desktop/web.log","topic":"web_log"},{"path":"/Users/caoyifan/Desktop/nazha.log","topic":"nazha_log"}]`
	_,err = cli.Put(ctx,"collect_log_conf",str)
	if err!= nil {
		fmt.Printf("put to etcd failed,err:%v",err)
		return
	}
	cancel()

	ctx,cancel = context.WithTimeout(context.Background(),time.Second)
	gr,err := cli.Get(ctx,"collect_log_conf")
	if err!= nil {
		fmt.Printf("get from etcd failed,err:%v",err)
		return
	}
	for _,ev := range gr.Kvs{
		fmt.Printf("key:%s value:%s\n",ev.Key,ev.Value)
	}
	cancel()
}