### 用Go语言实现的日志收集

#### 基于etcd+kafka+sarama实现

大致流程如下：

![image](https://github.com/elssm/Go-logagent/blob/main/5.png)

利用etcd来管理日志文件的path和topic

启动etcd，Linux命令行下执行etcd即可

```go
etcd
```
![image](https://github.com/elssm/Go-logagent/blob/main/4.png)

在etcd目录下执行

```go
go run main.go
```

#### 启动kafka

启动kafka之前先启动zookeeper

```shell
zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties &
```

之后启动kafka

```shell
kafka-server-start /usr/local/etc/kafka/server.properties &
```

#### logagent目录下的文件说明

##### - common：公共的日志结构

##### - conf：kafka和etcd的相关配置

##### - kafka：将日志发送到kafka中

##### - tailfile：从kafka中读取日志

启动logagent，在logagent目录下执行

```go
go build main.go
./main
```

#### 开启一个kafka消费者

```shell
sh kafka-console-consumer --bootstrap-server localhost:9092 --topic web_log --from-beginning
```

![image](https://github.com/elssm/Go-logagent/blob/main/1.png)
![image](https://github.com/elssm/Go-logagent/blob/main/2.png)
![image](https://github.com/elssm/Go-logagent/blob/main/3.png)
