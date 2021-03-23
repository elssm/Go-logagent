package kafka

import (
	//"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	Client sarama.SyncProducer
	MsgChan chan *sarama.ProducerMessage
)

//Init是初始化全局的kafka client
func Init(address []string,chanSize int64) (err error)  {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll //ACK
	config.Producer.Partitioner = sarama.NewRandomPartitioner //分区
	config.Producer.Return.Successes = true //确认

	//连接kafka
	Client,err = sarama.NewSyncProducer(address,config)
	if err != nil {
		logrus.Error("kafka:producer closed err :",err)
		return
	}
	MsgChan = make(chan *sarama.ProducerMessage,chanSize)
	//起一个后台的goroutine从msgchan中读数据
	go sendMsg()
	return
}

//从msgChan中读取消息发送给kafka
func sendMsg()  {
	for {
		select {
		case msg:= <- MsgChan:
			pid,offset,err := Client.SendMessage(msg)
			if err != nil {
				logrus.Warning("send msg failed,err:",err)
				return
			}
			logrus.Infof("send msg to kafka success. pid:%v offset:%v",pid,offset)
		}
	}
}