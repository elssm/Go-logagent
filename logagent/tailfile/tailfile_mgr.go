package tailfile

import (
	"github.com/sirupsen/logrus"
	"logagent/common"
)

//tailTask的管理者

type tailTaskMgr struct {
	tailTaskMap map[string]*tailTask
	collectEntryList []common.CollectEntry
	confChan chan []common.CollectEntry
}

var (
	ttMgr *tailTaskMgr
)

func Init(allConf []common.CollectEntry ) (err error)  {
	ttMgr = &tailTaskMgr{
		tailTaskMap: make(map[string]*tailTask,20),
		collectEntryList: allConf,
		confChan: make(chan []common.CollectEntry),
	}
	for _,conf := range allConf {
		tt := newTailTask(conf.Path,conf.Topic)
		err = tt.Init()
		if err != nil {
			logrus.Errorf("create tailObj for path:%s failed,err:%v\n",conf.Path,err)
			continue
		}
		logrus.Infof("create a tail task for path:%s success!\n",conf.Path)
		ttMgr.tailTaskMap[tt.path] = tt
		go tt.run()
	}
	go ttMgr.watch()

	return
}

func (t *tailTaskMgr)watch()  {
	for {
		newConf := <- t.confChan //取到值说明新的配置来了
		logrus.Infof("get new conf from etcd,conf:%v,start manage tailTask....",newConf)
		for _,conf := range newConf {
			//1.原来已经存在的任务就不用动
			if t.isExist(conf){
				continue
			}
			//2.原来没有的我要新创建一个tailTask任务
			tt := newTailTask(conf.Path,conf.Topic)
			err := tt.Init()
			if err != nil {
				logrus.Errorf("create tailObj for path:%s failed,err:%v\n",conf.Path,err)
				continue
			}
			logrus.Infof("create a tail task for path:%s success!\n",conf.Path)
			ttMgr.tailTaskMap[tt.path] = tt
			go tt.run()

		}
		//3.原来有的现在没有的要tailTask要停掉
		//找出tailTaskMap中存在，但是newConf中不存在的那些tailTask把它们都停掉
		for key,task := range t.tailTaskMap {
			var found bool
			for _,conf := range newConf {
				if key == conf.Path {
					found = true
					break
				}
			}
			if !found {
				//这个tailTask要停掉
				logrus.Infof("the task collect path:%s need to stop..",task.path)
				delete(t.tailTaskMap,key) //从管理类中删掉
				task.cancel()
			}

		}
	}

}

func (t *tailTaskMgr)isExist(conf common.CollectEntry) bool {
	_,ok := t.tailTaskMap[conf.Path]
	return ok
}

func SendNewConf(newConf []common.CollectEntry)  {
	ttMgr.confChan <- newConf
}