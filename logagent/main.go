package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/tailfile"
)

type Config struct {
	KafkaConfig `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig `ini:"etcd"`
}

type KafkaConfig struct {
	Address string `ini:"address"`
	//Topic string `ini:"topic"`
	ChanSize int64 `ini:"chan_size"`
}

type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

//真正的业务逻辑
//func run() (err error) {
//	for {
//		line,ok := <- tailfile.TailObj.Lines
//		if !ok {
//			logrus.Warn("tail file close reopen,filename:%s\n",tailfile.TailObj.Filename)
//			time.Sleep(time.Second) //读取出错等一秒
//			continue
//		}
//		//如果是空行就略过
//		if len(line.Text) == 0 {
//			continue
//		}
//		//利用通道将同步的代码改为异步的
//		//把读出来的一行日志包装成kafka里面的msg类型，丢到通道中
//		msg := &sarama.ProducerMessage{}
//		msg.Topic = "web_log"
//		msg.Value = sarama.StringEncoder(line.Text)
//		//丢到管道中
//		kafka.MsgChan <- msg
//	}
//}

func run()  {
	select {

	}
}

//日志收集的客户端
//收集指定目录下的日志文件，发送到kafka中
func main()  {
	var configObj = new(Config)
	//读配置文件
	//cfg,err := ini.Load("./conf/config.ini")
	//if err != nil {
	//	logrus.Errorf("load config failed,err:%v",err)
	//	return
	//}
	//kafkaAddr := cfg.Section("kafka").Key("address").String()
	//fmt.Println(kafkaAddr)

	err := ini.MapTo(configObj,"./conf/config.ini")
	if err != nil {
		logrus.Errorf("load config failed,err:%v",err)
		return
	}
	fmt.Printf("%#v",configObj)
	//初始化,连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address},configObj.KafkaConfig.ChanSize)
	if err != nil {
		logrus.Errorf("init kafka failed,err:%v",err)
	}
	logrus.Info("init kafka success!")

	//初始化etcd连接
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Errorf("init etcd failed err:%v",err)
		return
	}
	//从etcd中拉取要收集日志的配置项
	allConf,err := etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Errorf("get conf from etcd failed err:%v",err)
		return
	}
	fmt.Println(allConf)

	//监控etcd中configObj.EtcdConfig.CollectKey对应值的变化
	go etcd.WatchConf(configObj.EtcdConfig.CollectKey)

	//根据配置中的日志路径使用tail收集日志
	err = tailfile.Init(allConf) //把从etcd中获取的配置项传递到Init中
	if err != nil {
		logrus.Error("init tailfile failed,err:%v",err)
		return
	}
	logrus.Info("init tailfile success!")
	//把日志通过sarama发往kafka
	//TailObj --> log --> Client --> kafka
	//err = run()
	//if err != nil {
	//	logrus.Error("run failed,err :%v",err)
	//}
	run()
}