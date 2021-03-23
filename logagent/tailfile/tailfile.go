package tailfile

import (
	"context"
	"github.com/Shopify/sarama"
	"time"

	//"fmt"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logagent/common"
	"logagent/kafka"
)

type tailTask struct {
	path string
	topic string
	tObj *tail.Tail
	ctx context.Context
	cancel context.CancelFunc
}

//var (
//	TailObj *tail.Tail
//)
var (
	confChan chan []common.CollectEntry
)


//根据topic和path造一个tailTask对象
func newTailTask(path,topic string) *tailTask {
	ctx,cancel := context.WithCancel(context.Background())
	tt := &tailTask{
		path:path,
		topic: topic,
		ctx: ctx,
		cancel: cancel,
	}
	return tt
}

//使用tail包打开日志文件准备读
func (t *tailTask)Init() (err error)  {
	cfg := tail.Config{
		ReOpen:true,
		Follow:true,
		Location:&tail.SeekInfo{Offset:0,Whence:2},
		MustExist:false,
		Poll:true,
	}
	t.tObj,err = tail.TailFile(t.path,cfg)
	return
}

//实际读日志往kafka里发送的方法
func (t *tailTask)run()  {
	logrus.Infof("collect for path:%s is running... ",t.path)
	for {
		select {
		case <- t.ctx.Done(): //只要调用t.cancel()就回收到信号
			logrus.Infof("path:%s is stopping...",t.path)
			return
		case line,ok := <- t.tObj.Lines:
			if !ok {
				logrus.Warn("tail file close reopen,path:%s\n",t.path)
				time.Sleep(time.Second) //读取出错等一秒
				continue
			}
			//如果是空行就略过
			if len(line.Text) == 0 {
				continue
			}
			//利用通道将同步的代码改为异步的
			//把读出来的一行日志包装成kafka里面的msg类型，丢到通道中
			msg := &sarama.ProducerMessage{}
			msg.Topic = t.topic
			msg.Value = sarama.StringEncoder(line.Text)
			//丢到管道中
			kafka.MsgChan <- msg

		}
		//line,ok := <- t.tObj.Lines
		//if !ok {
		//	logrus.Warn("tail file close reopen,path:%s\n",t.path)
		//	time.Sleep(time.Second) //读取出错等一秒
		//	continue
		//}
		////如果是空行就略过
		//if len(line.Text) == 0 {
		//	continue
		//}
		////利用通道将同步的代码改为异步的
		////把读出来的一行日志包装成kafka里面的msg类型，丢到通道中
		//msg := &sarama.ProducerMessage{}
		//msg.Topic = t.topic
		//msg.Value = sarama.StringEncoder(line.Text)
		////丢到管道中
		//kafka.MsgChan <- msg
	}
}

//初始化tailTask，为每一个日志文件造一个单独的tailTask
//func Init(allConf []common.CollectEntry ) (err error)  {
//	//allConf里面存了若干个日志的收集项
//	//cfg := tail.Config{
//	//	ReOpen:true,
//	//	Follow:true,
//	//	Location:&tail.SeekInfo{Offset:0,Whence:2},
//	//	MustExist:false,
//	//	Poll:true,
//	//}
//	for _,conf := range allConf {
//		tt := newTailTask(conf.Path,conf.Topic)
//		err = tt.Init()
//		if err != nil {
//			logrus.Errorf("create tailObj for path:%s failed,err:%v\n",conf.Path,err)
//			continue
//		}
//		logrus.Infof("create a tail task for path:%s success!\n",conf.Path)
//		go tt.run()
//	}
//
//	//初始化新配置的管道
//	confChan = make(chan []common.CollectEntry) //做一个阻塞channel
//	newConf := <- confChan //取到值说明新的配置来了
//	logrus.Infof("get new conf from etcd,conf:%v",newConf)
//	return
//}

//func SendNewConf(newConf []common.CollectEntry)  {
//	confChan <- newConf
//}