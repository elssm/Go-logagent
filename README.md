### 用Go语言实现的日志收集

#### 基于etcd+kafka+sarama实现

利用etcd来管理日志文件的path和topic

启动etcd，在etcd目录下执行

```go
go run main.go
```

#### logagent目录下的文件说明

##### - common：公共的日志结构

##### - conf：kafka和etcd的相关配置

##### - kafka：将日志发送到kafka中

##### - tailfile：从kafka中读取日志

启动logagent，在logagent目录下执行

```go
go run main.go
```

