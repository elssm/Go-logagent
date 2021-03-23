package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/tailfile"
	"time"
	"go.etcd.io/etcd/clientv3"
)

//type collectEntry struct {
//	Path string `json:"path"`
//	Topic string `json:"topic"`
//}

var (
	client *clientv3.Client
)

//etcd相关操作

func Init(address []string) (err error)  {
	client,err = clientv3.New(clientv3.Config{
		Endpoints: address,
		DialTimeout: time.Second*5,
	})
	if err!=nil {
		fmt.Printf("connect to etcd failed , err: %v",err)
		return
	}
	return
}

func GetConf(key string) (collectEntryList []common.CollectEntry,err error)  {
	ctx,cancel := context.WithTimeout(context.Background(),time.Second*2)
	defer cancel()
	resp,err := client.Get(ctx,key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key:%s failed,err:%v",key,err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warnf("get len:0 conf from etcd by key:%s",key)
		return
	}
	ret := resp.Kvs[0]
	fmt.Println(ret.Value)
	//ret.Value //json格式字符串
	err = json.Unmarshal(ret.Value,&collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed,err:%v",err)
		return
	}
	return
}

//监控etcd中日志收集配置项变化的函数
func WatchConf(key string)  {
	for {
		watchCh := client.Watch(context.Background(), key)
		var newConf []common.CollectEntry
		for wresp := range watchCh {
			logrus.Info("get new conf from etcd!")
			for _, evt := range wresp.Events {
				fmt.Printf("type:%s key:%s value:%s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshal new conf failed,err:%v", err)
					continue
				}
				//告诉tailfile模块应该启用新的配置
				tailfile.SendNewConf(newConf) //没有人接收就是阻塞的
			}
		}
	}
}