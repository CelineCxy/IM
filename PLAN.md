# Go 后端开发 40 天系统学习计划

> 适合人群：已经学习过 Go 基础语法，希望系统学习 Go 后端、并发、数据库、中间件、微服务、分布式和工程化开发。
> 学习周期：40 天
> 每日建议时间：3～5 小时
> 最终目标：独立完成一个 Go 分布式 IM 系统，并具备 Go 后端面试所需的核心知识体系。

---

# 一、整体学习路线

```text
Go 工程化
    ↓
Go 并发编程
    ↓
网络编程
    ↓
Gin Web 开发
    ↓
MySQL + Redis
    ↓
WebSocket
    ↓
Protobuf + gRPC
    ↓
Etcd + 服务发现
    ↓
Kafka 消息队列
    ↓
分布式消息系统
    ↓
Docker + Nginx
    ↓
日志 + 监控 + 链路追踪
    ↓
pprof 性能优化
    ↓
分布式 IM 综合项目
```

最终项目技术栈：

```text
Go
Gin
WebSocket
gRPC
Protobuf
MySQL
Redis
MongoDB
Etcd
Kafka
Docker
Docker Compose
Nginx
Zap
Viper
Prometheus
Grafana
OpenTelemetry
Jaeger
pprof
```

---

# 二、学习项目目录建议

整个 40 天只维护一个代码仓库：

```text
go-backend-learning/
├── day01/
├── day02/
├── day03/
├── ...
├── day40/
│
├── projects/
│   ├── gin-user-service/
│   ├── websocket-chat/
│   ├── grpc-demo/
│   ├── kafka-demo/
│   └── distributed-im/
│
├── notes/
│   ├── golang/
│   ├── mysql/
│   ├── redis/
│   ├── grpc/
│   ├── kafka/
│   └── interview/
│
└── README.md
```

每天至少保证：

```text
1. 学一个核心知识模块
2. 写一个 Demo
3. 完成一个小功能
4. 总结 3～5 个面试问题
5. 提交一次 Git
```

---

# 第一阶段：Go 工程能力与并发

## Day 1～Day 8

---

# Day 1：Go 项目工程化

## 今日目标

理解一个正规的 Go 后端项目应该如何组织，而不是所有代码全部放在 `main.go`。

## 学习内容

### 1. Go Module

重点掌握：

```bash
go mod init
go mod tidy
go mod download
go get
go list
```

理解：

```text
go.mod
go.sum
```

分别有什么作用。

### 2. Go package

理解：

```text
package main
package service
package repository
package model
```

重点理解：

* 一个目录为什么通常对应一个 package
* 大写方法为什么可以跨包访问
* 小写方法为什么只能包内访问
* init() 的执行时机
* import cycle 为什么会报错

### 3. Go 项目目录

创建：

```text
user-service/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
├── pkg/
├── configs/
├── scripts/
├── go.mod
└── README.md
```

理解：

```text
cmd
internal
pkg
configs
```

分别解决什么问题。

## 动手任务

实现：

```text
main
 ↓
handler
 ↓
service
 ↓
repository
```

先不用数据库。

例如：

```text
GET /user?id=1
```

repository 返回模拟用户：

```json
{
  "id": 1,
  "name": "Tom"
}
```

## 今日必须搞懂

* `internal` 有什么特殊性？
* `pkg` 和 `internal` 有什么区别？
* Go 为什么不推荐特别复杂的目录？
* Go Module 解决了什么问题？

## 今日产出

* [X] 创建 Go 学习仓库
* [X] 创建规范目录
* [X] 完成简单分层 Demo
* [X] 上传 Git

---

# Day 2：Interface、Error 与泛型

## 今日目标

掌握 Go 工程代码里最频繁出现的抽象方式。

## 学习内容

### 1. Interface

重点理解：

```go
type UserRepository interface {
    GetUser(id int64) (*User, error)
}
```

然后实现：

```go
type MysqlUserRepository struct {
}
```

思考：

```text
为什么 Service 依赖 interface，而不是直接依赖 MysqlUserRepository？
```

理解依赖倒置思想。

### 2. 类型断言

练习：

```go
v, ok := x.(string)
```

以及：

```go
switch v := x.(type) {

}
```

### 3. Error

学习：

```go
errors.New()
fmt.Errorf()
errors.Is()
errors.As()
errors.Join()
```

自己定义：

```go
var ErrUserNotFound = errors.New("user not found")
```

Service 返回：

```go
fmt.Errorf("get user failed: %w", ErrUserNotFound)
```

### 4. 泛型

学习：

```go
func Max[T constraints.Ordered](a, b T) T
```

重点理解：

```text
泛型解决什么问题？
什么时候不要使用泛型？
```

## 动手任务

设计统一返回结构：

```go
type Response[T any] struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    T      `json:"data"`
}
```

分别返回：

```text
Response[User]
Response[[]User]
```

## 面试复盘

回答：

1. interface 和 Java interface 有什么区别？
2. Go interface 是隐式实现是什么意思？
3. `errors.Is` 和 `==` 有什么区别？
4. `errors.As` 是干什么的？
5. Go 泛型适合什么场景？

---

# Day 3：Context

## 今日目标

真正理解 Go 后端为什么到处都是：

```go
ctx context.Context
```

## 学习内容

学习：

```go
context.Background()
context.TODO()
context.WithCancel()
context.WithTimeout()
context.WithDeadline()
context.WithValue()
```

理解 Context 的核心作用：

```text
超时控制
取消传播
请求生命周期管理
请求级数据传递
```

## 动手任务 1

模拟一个数据库请求：

```text
请求耗时：5 秒
超时时间：2 秒
```

要求：

```text
2 秒后自动取消请求。
```

## 动手任务 2

创建：

```text
HTTP Request
    ↓
Service
    ↓
Repository
```

Context 一直向下传递：

```go
Handler(ctx)
    ↓
Service(ctx)
    ↓
Repository(ctx)
```

## 动手任务 3

模拟两个 goroutine：

```text
goroutine A
goroutine B
```

只要主请求取消：

```text
A、B 同时退出。
```

## 今日重点

搞懂：

```text
为什么 Context 要作为第一个参数？
为什么不应该把 Context 存在 struct 中？
为什么不能滥用 WithValue？
```

---

# Day 4：Goroutine 与 GMP

## 今日目标

开始进入 Go 最重要的知识模块：

```text
并发编程
```

## 理论学习

掌握：

```text
Process
Thread
Coroutine
Goroutine
```

理解：

```text
GMP 调度模型
```

其中：

```text
G = Goroutine
M = Machine Thread
P = Processor
```

理解：

```text
G 为什么比线程轻量？
goroutine 初始栈为什么很小？
goroutine 如何扩容？
```

了解：

```text
work stealing
hand off
```

## 动手任务

普通版本：

```text
下载 10 个 URL
```

改造成：

```text
10 个 goroutine 并发下载。
```

记录：

```text
串行耗时
并发耗时
```

## 第二个任务

创建：

```text
100000 个 goroutine
```

观察：

```text
内存变化
goroutine 数量
```

使用：

```go
runtime.NumGoroutine()
```

## 面试题

必须能够回答：

1. goroutine 和线程有什么区别？
2. GMP 是什么？
3. P 的数量由什么决定？
4. GOMAXPROCS 是什么？
5. goroutine 会不会泄漏？

---

# Day 5：Channel

## 今日目标

掌握 Go 最核心的并发通信机制。

## 学习内容

学习：

```go
ch := make(chan int)
```

以及：

```go
ch := make(chan int, 10)
```

理解：

```text
无缓冲 Channel
有缓冲 Channel
```

学习：

```go
ch <- value
value := <-ch
close(ch)
```

以及：

```go
for value := range ch
```

学习：

```go
select
```

## 动手任务 1

实现：

```text
Producer
    ↓
Channel
    ↓
Consumer
```

生产 100 个任务。

## 动手任务 2

实现一个 Worker Pool：

```text
任务队列
   ↓
Channel
   ↓
5 个 Worker
   ↓
并发执行
```

任务：

```text
计算 1～1000 每个数字的平方。
```

## 动手任务 3

加入：

```go
context.WithTimeout()
```

如果任务执行超过 5 秒：

```text
所有 Worker 停止。
```

## 必须理解

```text
关闭 Channel 后还能读取吗？
关闭 Channel 后还能写入吗？
向 nil channel 写数据会怎样？
select default 有什么作用？
```

---

# Day 6：Mutex、RWMutex、Atomic

## 今日目标

学习 Go 共享内存并发控制。

## 学习内容

学习：

```go
sync.Mutex
sync.RWMutex
sync.WaitGroup
sync.Once
sync.Cond
sync.Map
```

以及：

```go
sync/atomic
```

## 动手任务 1

创建：

```go
count := 0
```

启动：

```text
1000 个 goroutine
```

每个执行：

```go
count++
```

先观察错误结果。

然后分别使用：

```text
Mutex
Atomic
```

实现正确版本。

## 动手任务 2

自己实现：

```go
type SafeMap struct {
    mu sync.RWMutex
    data map[string]string
}
```

实现：

```text
Get
Set
Delete
```

## 动手任务 3

使用：

```bash
go run -race main.go
```

观察 Race Detector 如何发现数据竞争。

## 面试题

回答：

```text
Mutex 和 RWMutex 区别？
Mutex 是否可重入？
channel 和 mutex 应该怎么选？
CAS 是什么？
atomic 为什么比 mutex 快？
```

---

# Day 7：高级并发模型

## 今日目标

掌握真正会在业务中使用的并发模式。

## 学习内容

学习：

```text
Worker Pool
Pipeline
Fan-out
Fan-in
Semaphore
限流
```

## 动手任务

实现图片处理流水线：

```text
读取图片
   ↓
并发压缩
   ↓
并发上传
   ↓
结果汇总
```

分别设计：

```text
Downloader Pool
Processor Pool
Uploader Pool
```

每一层通过 Channel 通信。

## 第二个任务

实现最大并发数：

```text
最多同时执行 10 个任务。
```

可以使用：

```go
make(chan struct{}, 10)
```

实现 Semaphore。

## 重点问题

分析：

```text
如果某个阶段特别慢，会发生什么？
Channel buffer 应该设置多大？
如何避免 goroutine 泄漏？
```

---

# Day 8：Go Test、Benchmark 与 Race

## 今日目标

掌握 Go 工程项目测试能力。

## 学习内容

学习：

```bash
go test
go test -v
go test ./...
go test -race
go test -bench=.
```

学习：

```go
func TestXXX(t *testing.T)
```

和：

```go
func BenchmarkXXX(b *testing.B)
```

## 动手任务

给 SafeMap 编写：

```text
单元测试
并发测试
Benchmark
```

比较：

```text
Mutex
RWMutex
sync.Map
```

在：

```text
读多写少
读写均衡
```

场景下的性能。

## 今日产出

写笔记：

```text
什么时候使用 sync.Map？
RWMutex 一定比 Mutex 快吗？
Benchmark 怎么避免测试误差？
```

---

# 第二阶段：网络与 Web 后端

# Day 9～Day 13

---

# Day 9：TCP 网络编程

## 学习内容

理解：

```text
TCP 三次握手
TCP 四次挥手
长连接
短连接
粘包
拆包
KeepAlive
```

学习 Go：

```go
net.Listen()
listener.Accept()
net.Dial()
conn.Read()
conn.Write()
```

## 动手任务

实现：

```text
TCP Server
TCP Client
```

支持：

```text
多个客户端同时连接。
```

客户端发送：

```text
hello
```

Server 返回：

```text
echo: hello
```

## 进阶任务

设计消息协议：

```text
4 字节消息长度
+
消息 Body
```

自己解决 TCP 粘包问题。

---

# Day 10：HTTP 与 net/http

## 学习目标

不要直接上 Gin，先理解 Gin 底层依赖的 HTTP Server。

## 学习内容

理解：

```text
HTTP Request
HTTP Response
Header
Body
Cookie
Session
Keep-Alive
HTTP 状态码
RESTful
```

学习：

```go
http.HandleFunc()
http.ListenAndServe()
http.Request
http.ResponseWriter
```

## 动手任务

不用 Gin 实现：

```text
POST /users
GET /users/{id}
DELETE /users/{id}
```

实现：

```text
JSON 解析
JSON 返回
参数校验
错误处理
```

---

# Day 11：Gin

## 学习内容

学习：

```text
Router
RouterGroup
Middleware
Query
Param
JSON Binding
Recovery
Logger
```

## 动手项目

创建：

```text
User Service
```

接口：

```text
POST   /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
```

统一返回：

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

## 必须实现

中间件：

```text
RequestID
日志
耗时统计
Recovery
```

---

# Day 12：标准 Web 分层架构

## 今日目标

把昨天代码重构成真正后端项目。

目录：

```text
internal/
├── handler/
├── service/
├── repository/
├── model/
├── middleware/
└── router/
```

调用链：

```text
HTTP
 ↓
Handler
 ↓
Service
 ↓
Repository
 ↓
Database
```

## 动手任务

实现：

```text
UserHandler
UserService
UserRepository
```

Repository 暂时使用：

```text
Map 模拟数据库。
```

## 重点理解

为什么：

```text
Handler 不直接访问数据库？
Service 为什么不应该知道 HTTP？
Repository 为什么只负责数据访问？
```

---

# Day 13：JWT 与认证授权

## 学习内容

理解：

```text
Authentication
Authorization
```

学习 JWT：

```text
Header
Payload
Signature
```

实现：

```text
Access Token
Refresh Token
```

## 动手任务

接口：

```text
POST /register
POST /login
POST /refresh
GET /profile
```

其中：

```text
/profile
```

必须通过 JWT。

## 实现中间件

```text
Authorization: Bearer token
```

解析：

```text
user_id
username
expire_time
```

## 进阶

简单实现 RBAC：

```text
admin
user
```

限制：

```text
/admin/*
```

只有 admin 可访问。

---

# 第三阶段：数据库与缓存

# Day 14～Day 19

---

# Day 14：MySQL 核心

## 学习内容

重点复习：

```text
B+ Tree
聚簇索引
二级索引
联合索引
最左匹配
覆盖索引
回表
索引下推
```

学习：

```sql
EXPLAIN
```

重点看：

```text
type
key
rows
Extra
```

## 动手任务

创建：

```sql
user
conversation
message
```

三张表。

设计：

```text
主键
普通索引
联合索引
唯一索引
```

## SQL 练习

查询：

```text
某用户最近 20 个会话
某会话最近 50 条消息
某用户某时间段消息
```

使用：

```sql
EXPLAIN
```

分析是否走索引。

---

# Day 15：Go 操作 MySQL

## 学习

使用：

```go
database/sql
```

学习连接池参数：

```go
SetMaxOpenConns()
SetMaxIdleConns()
SetConnMaxLifetime()
```

## 动手任务

实现：

```text
CreateUser
GetUser
UpdateUser
DeleteUser
```

## 必须实现事务

模拟：

```text
创建订单
 ↓
扣减库存
 ↓
生成支付记录
```

任何一步失败：

```text
Rollback
```

全部成功：

```text
Commit
```

## 面试题

回答：

```text
数据库连接池为什么重要？
连接数是不是越多越好？
事务 ACID 是什么？
MySQL MVCC 是什么？
```

---

# Day 16：GORM

## 学习内容

学习：

```text
Model
Create
First
Find
Updates
Delete
Transaction
Raw SQL
Hooks
```

## 动手任务

把 Day15 Repository 改成：

```text
GORM
```

定义：

```go
type User struct {}
type Conversation struct {}
type Message struct {}
```

## 重点

不要只会：

```go
db.Find()
```

还需要理解：

```text
ORM 会产生什么 SQL？
什么时候应该直接写 SQL？
N+1 查询是什么？
```

---

# Day 17：Redis

## 学习内容

完整学习：

```text
String
Hash
List
Set
ZSet
Bitmap
HyperLogLog
```

重点掌握：

```text
TTL
Expire
Pipeline
Transaction
Lua
```

## 动手任务

使用：

```text
go-redis
```

实现：

```text
用户缓存
验证码
在线状态
排行榜
```

例如：

```text
user:1
online:10001
sms:138xxxx
ranking:user
```

---

# Day 18：缓存与分布式锁

## 学习内容

重点理解：

```text
缓存穿透
缓存击穿
缓存雪崩
缓存一致性
```

学习：

```text
Cache Aside
```

流程：

```text
读缓存
 ↓
缓存没有
 ↓
读数据库
 ↓
写缓存
```

## 动手任务 1

实现用户查询缓存：

```text
GET /users/:id
```

增加 Redis。

## 动手任务 2

实现：

```text
SET NX EX
```

分布式锁。

## 进阶

使用 Lua：

```text
判断 value
+
删除锁
```

保证释放锁原子性。

---

# Day 19：MongoDB

## 今日目标

不要求特别深入，但要知道什么时候 MongoDB 更适合。

## 学习内容

理解：

```text
Document
Collection
ObjectId
Index
Aggregation
```

## 动手任务

存储：

```json
{
  "conversation_id": 10001,
  "sender": 1,
  "receiver": 2,
  "content": "hello",
  "timestamp": 1720000000
}
```

实现：

```text
保存聊天记录
查询聊天记录
分页查询
```

最后比较：

```text
MySQL 保存消息
MongoDB 保存消息
```

各自优缺点。

---

# 第四阶段：IM 与微服务

# Day 20～Day 30

---

# Day 20：WebSocket

## 学习内容

理解：

```text
HTTP Upgrade
WebSocket Handshake
Full Duplex
Ping
Pong
Close
```

## 动手任务

使用 Go WebSocket 库实现：

```text
Client
 ↓
WebSocket
 ↓
Server
```

支持：

```text
连接
发消息
接收消息
断开连接
```

---

# Day 21：WebSocket IM Gateway

## 今日目标

把 WebSocket Demo 升级成真正的连接网关。

## 数据结构

设计：

```go
type Client struct {
    UserID int64
    Conn   *websocket.Conn
    Send   chan []byte
}
```

以及：

```go
type Hub struct {
    clients map[int64]*Client
}
```

## 实现

```text
用户上线
用户下线
发送消息
连接关闭
异常断开
心跳检测
```

## 架构

```text
Read Goroutine
       ↓
 WebSocket
       ↑
Write Goroutine
```

读写协程解耦。

## 重点问题

思考：

```text
为什么不能多个 goroutine 同时 Write WebSocket？
慢客户端怎么办？
Send Channel 满了怎么办？
```

---

# Day 22：Protobuf

## 学习内容

安装：

```text
protoc
protoc-gen-go
protoc-gen-go-grpc
```

学习：

```protobuf
message
enum
repeated
service
rpc
```

## 动手任务

定义：

```protobuf
message SendMessageRequest {
    int64 sender_id = 1;
    int64 receiver_id = 2;
    string content = 3;
}
```

以及：

```protobuf
service MessageService {
    rpc SendMessage(SendMessageRequest)
        returns (SendMessageResponse);
}
```

执行代码生成。

---

# Day 23：gRPC

## 学习内容

理解：

```text
Unary RPC
Server Streaming
Client Streaming
Bidirectional Streaming
```

## 动手任务

拆成：

```text
Gateway Service
Message Service
```

Gateway 调：

```text
MessageService.SendMessage()
```

## 第二个任务

实现一个：

```text
Server Streaming
```

例如：

```text
历史消息流式返回。
```

---

# Day 24：gRPC 工程化

## 学习内容

学习：

```text
Interceptor
Metadata
Deadline
Status Code
TLS
KeepAlive
```

## 动手任务

写：

```text
AuthInterceptor
LoggingInterceptor
RecoveryInterceptor
TraceInterceptor
```

把：

```text
user_id
trace_id
```

放入 Metadata。

## 超时

Gateway 设置：

```text
500ms
```

如果 Message Service 超时：

```text
主动取消。
```

---

# Day 25：微服务拆分

## 今日目标

正式开始做：

```text
Distributed IM
```

拆成：

```text
Gateway Service
User Service
Message Service
```

职责：

### Gateway

```text
WebSocket
长连接
用户连接
消息转发
```

### User Service

```text
注册
登录
用户信息
好友
```

### Message Service

```text
消息存储
消息查询
离线消息
```

## 今日任务

画架构图：

```text
Client
   ↓
Gateway
 ┌───────┴────────┐
 ↓                ↓
User Service   Message Service
```

完成服务间 gRPC 调用。

---

# Day 26：Etcd 服务注册发现

## 学习内容

理解：

```text
KV
Lease
TTL
Watch
Revision
```

## 动手任务

服务启动：

```text
Message Service
```

自动注册：

```text
/services/message/instance-1
```

Value：

```text
10.0.0.1:9001
```

通过 Lease：

```text
10 秒 TTL
```

持续 KeepAlive。

## Gateway

通过：

```text
Watch
```

实时获取节点变化。

实现：

```text
服务上线
服务下线
故障摘除
```

---

# Day 27：负载均衡与一致性哈希

## 学习

理解：

```text
Round Robin
Random
Least Connections
Hash
Consistent Hash
```

## 动手任务

自己实现：

```go
type HashRing struct {}
```

实现：

```text
AddNode
RemoveNode
GetNode
```

创建：

```text
node1
node2
node3
```

模拟：

```text
10000 个 userID
```

观察节点分布。

然后删除 node2：

```text
观察多少用户发生迁移。
```

---

# Day 28：Kafka

## 学习内容

掌握：

```text
Broker
Topic
Partition
Producer
Consumer
Consumer Group
Offset
ISR
Replica
```

重点理解：

```text
Kafka 为什么吞吐量高？
Partition 为什么影响顺序？
Consumer Group 有什么作用？
```

## 动手任务

创建：

```text
Topic: im-message
```

Gateway：

```text
Producer
```

Message Service：

```text
Consumer
```

消息：

```json
{
  "message_id": 100001,
  "sender_id": 1,
  "receiver_id": 2,
  "content": "hello"
}
```

---

# Day 29：消息可靠性

## 今日重点

研究：

```text
消息丢失
重复消息
乱序消息
消息积压
```

## 动手任务

设计：

```text
MessageID
```

例如：

```text
Snowflake / UUID
```

消费者保存消息前：

```text
检查 MessageID 是否已经处理。
```

实现：

```text
幂等消费
```

## 重试

第一次失败：

```text
retry 1
retry 2
retry 3
```

仍失败：

```text
进入 Dead Letter Queue。
```

---

# Day 30：IM 消息可靠投递

## 今日目标

完成完整消息链路。

```text
Client A
   ↓
Gateway A
   ↓
Kafka
   ↓
Message Service
   ↓
MySQL
   ↓
Gateway B
   ↓
Client B
```

## 消息状态

定义：

```text
Created
Sent
Delivered
Read
```

## 实现 ACK

Client 收到消息后：

```json
{
  "type": "ack",
  "message_id": 10001
}
```

Server 更新：

```text
Delivered
```

## 离线消息

如果用户不在线：

```text
消息保存数据库
```

重新登录：

```text
拉取未读消息。
```

---

# 第五阶段：工程化

# Day 31～Day 37

---

# Day 31：Docker

## 学习

掌握：

```text
Image
Container
Dockerfile
Volume
Network
Registry
```

## 动手任务

为：

```text
Gateway
User Service
Message Service
```

分别写 Dockerfile。

使用：

```text
multi-stage build
```

例如：

```dockerfile
FROM golang AS builder

FROM alpine
```

减少镜像大小。

---

# Day 32：Docker Compose + Nginx

## 动手任务

Compose 启动：

```text
MySQL
Redis
Etcd
Kafka
Gateway
User Service
Message Service
```

完成：

```bash
docker compose up -d
```

整个系统运行。

## Nginx

配置：

```text
/api/*
```

转发 HTTP。

```text
/ws
```

转发 WebSocket。

理解：

```text
反向代理
负载均衡
```

---

# Day 33：Viper + Zap

## 配置管理

使用 Viper 管理：

```yaml
server:
  port: 8080

mysql:
  host: mysql
  port: 3306

redis:
  host: redis
```

区分：

```text
dev
test
prod
```

## 日志

使用 Zap：

```text
Debug
Info
Warn
Error
```

输出 JSON：

```json
{
  "level": "info",
  "trace_id": "xxx",
  "user_id": 10001,
  "message": "send message"
}
```

---

# Day 34：Prometheus + Grafana

## 学习指标

理解：

```text
Counter
Gauge
Histogram
Summary
```

## 动手任务

增加：

```text
HTTP QPS
WebSocket 在线人数
消息发送数量
消息失败数量
gRPC 请求耗时
```

例如：

```text
im_online_users
im_message_total
grpc_request_duration
```

暴露：

```text
/metrics
```

Prometheus 拉取。

Grafana 创建 Dashboard。

---

# Day 35：OpenTelemetry + Jaeger

## 今日目标

实现分布式链路追踪。

请求：

```text
HTTP
 ↓
Gateway
 ↓
gRPC
 ↓
Message Service
 ↓
MySQL
```

每一层创建：

```text
Span
```

传递：

```text
TraceID
```

最终可以在 Jaeger 中看到完整调用链。

重点理解：

```text
Trace
Span
Parent Span
Context Propagation
```

---

# Day 36：pprof 性能分析

## 学习

使用：

```text
net/http/pprof
```

查看：

```text
CPU Profile
Heap
Goroutine
Mutex
Block
```

命令：

```bash
go tool pprof
```

## 动手任务

故意写一个：

```text
CPU 消耗高程序
内存泄漏程序
goroutine 泄漏程序
```

分别使用 pprof 找出问题。

## 学习 GC

理解：

```text
Heap
Stack
GC
STW
三色标记
写屏障
```

了解：

```bash
go build -gcflags="-m"
```

观察逃逸分析。

---

# Day 37：高并发优化

## 今日学习

学习：

```text
连接池
对象池
sync.Pool
批量操作
减少锁竞争
限流
超时
熔断
降级
```

## 动手任务

针对 IM Gateway 做压测。

记录：

```text
连接数
QPS
P99
CPU
内存
goroutine 数量
```

优化：

```text
减少对象创建
减少锁粒度
批量写数据库
Redis Pipeline
Kafka Batch
```

优化前后做对比。

---

# 第六阶段：综合项目与面试

# Day 38～Day 40

---

# Day 38：完整 IM 系统整合

今天不要学习新技术。

只做整合。

完整架构：

```text
                 ┌──────────────┐
                 │    Client    │
                 └──────┬───────┘
                        │
                   WebSocket
                        │
                 ┌──────▼───────┐
                 │    Nginx     │
                 └──────┬───────┘
                        │
                 ┌──────▼───────┐
                 │ IM Gateway   │
                 └───┬──────┬───┘
                     │      │
                   gRPC    Redis
                     │
              ┌──────▼──────┐
              │Message Svc  │
              └──────┬──────┘
                     │
                   Kafka
                     │
             ┌───────▼───────┐
             │ Message Worker│
             └───────┬───────┘
                     │
                   MySQL
```

Etcd：

```text
服务注册发现
```

Prometheus：

```text
指标采集
```

Jaeger：

```text
链路追踪
```

完成核心链路测试。

---

# Day 39：异常场景测试

今天专门测试异常。

## 场景 1

```text
Redis 挂掉
```

观察：

```text
系统是否还能运行？
```

## 场景 2

```text
Message Service 挂掉
```

观察：

```text
Etcd 是否摘除节点？
```

## 场景 3

```text
Kafka Consumer 挂掉
```

恢复后：

```text
消息是否继续消费？
```

## 场景 4

```text
用户断网
```

测试：

```text
心跳超时
连接清理
重新连接
离线消息
```

## 场景 5

重复消息：

```text
同一个 MessageID 发送两次
```

确认：

```text
不会重复入库。
```

---

# Day 40：项目总结 + 面试复盘

## 第一件事：完善 README

包含：

```text
项目背景
系统架构
技术栈
数据库设计
Redis 设计
消息链路
服务发现
消息可靠性
性能测试
部署方式
```

---

## 第二件事：画架构图

需要能够自己讲清楚：

```text
客户端
 ↓
Gateway
 ↓
gRPC
 ↓
Kafka
 ↓
Message Service
 ↓
MySQL
```

Redis 和 Etcd 分别做什么。

---

## 第三件事：准备项目介绍

能够完整回答：

```text
为什么做这个项目？
系统整体架构是什么？
为什么使用 WebSocket？
为什么使用 gRPC？
为什么使用 Kafka？
Redis 用在哪里？
Etcd 用在哪里？
为什么使用一致性哈希？
消息如何保证不丢？
消息如何避免重复？
用户掉线怎么办？
Gateway 挂了怎么办？
Message Service 挂了怎么办？
```

---

## 第四件事：准备 Go 高频面试题

至少整理以下模块：

### Go

```text
GMP
Channel
Mutex
Context
Interface
Map
Slice
GC
逃逸分析
pprof
```

### MySQL

```text
索引
B+Tree
MVCC
事务
锁
慢 SQL
```

### Redis

```text
数据结构
缓存问题
持久化
分布式锁
高可用
```

### 网络

```text
TCP
HTTP
WebSocket
```

### 微服务

```text
gRPC
Etcd
服务注册发现
负载均衡
超时
重试
```

### MQ

```text
Kafka
顺序消息
重复消费
消息积压
可靠性
```

---

# 三、每天统一执行模板

以后每天学习时，都按照下面模板记录：

````markdown
# Day XX 学习记录

## 1. 今日学习目标

- [ ] 目标 1
- [ ] 目标 2
- [ ] 目标 3

## 2. 理论知识

### 知识点 1

自己的理解：

### 知识点 2

自己的理解：

---

## 3. 今日代码

项目路径：

```text
projects/xxx
````

实现功能：

* [ ] 功能 1
* [ ] 功能 2

---

## 4. 遇到的问题

### 问题 1

错误：

```text
xxx
```

原因：

```text
xxx
```

解决：

```text
xxx
```

---

## 5. 今日面试题

### Q1：xxx？

回答：

xxx

### Q2：xxx？

回答：

xxx

---

## 6. 今日总结

今天主要掌握：

1. xxx
2. xxx
3. xxx

还有问题：

1. xxx
2. xxx

````

---

# 四、40 天学习优先级

不是所有技术都需要同样深入。

## S 级：必须深入

```text
Go 并发
GMP
Channel
Mutex
Context
MySQL
Redis
TCP
HTTP
WebSocket
gRPC
````

这些属于 Go 后端核心。

---

## A 级：项目需要掌握

```text
Gin
GORM
Protobuf
Etcd
Kafka
Docker
Nginx
```

至少达到：

```text
知道原理
能自己使用
能解释为什么使用
```

---

## B 级：了解 + 能使用

```text
MongoDB
Prometheus
Grafana
OpenTelemetry
Jaeger
```

不需要一开始特别深入。

---

# 五、每天学习时间建议

如果每天有 4 小时：

```text
19:00 - 20:00
理论学习

20:00 - 22:00
代码实践

22:00 - 22:30
整理笔记

22:30 - 23:00
面试题复盘
```

比例建议：

```text
理论：30%
编码：50%
总结：20%
```

不要变成：

```text
看视频：90%
写代码：10%
```

正确方式应该是：

```text
学一个知识点
    ↓
马上写 Demo
    ↓
制造一个 Bug
    ↓
解决 Bug
    ↓
总结原因
    ↓
整理成面试回答
```

---

# 六、最终验收标准

40 天完成以后，不要求成为 Go 专家，但至少需要达到以下标准。

## Go

可以解释：

```text
GMP
Goroutine
Channel
Mutex
Context
GC
Interface
逃逸分析
pprof
```

---

## 后端

可以独立写：

```text
Gin REST API
JWT 登录
MySQL CRUD
Redis Cache
WebSocket
gRPC
Kafka Producer / Consumer
Etcd 服务注册发现
```

---

## 分布式

能够解释：

```text
服务发现
负载均衡
一致性哈希
分布式锁
消息幂等
消息重试
离线消息
消息 ACK
```

---

## 工程

可以独立完成：

```text
Docker 部署
Docker Compose 编排
Nginx 反向代理
结构化日志
Prometheus 监控
Jaeger 链路追踪
pprof 性能排查
```

---

# 七、最终项目

项目名称可以定义为：

```text
Distributed IM
```

或者：

```text
GoIM
```

最终项目技术架构：

```text
Go + Gin
    +
WebSocket
    +
gRPC / Protobuf
    +
Redis
    +
MySQL
    +
Kafka
    +
Etcd
    +
Docker
    +
Nginx
    +
Prometheus
    +
Grafana
    +
OpenTelemetry
    +
Jaeger
```

完成之后，这个项目应该不仅做到：

> “项目能运行。”

而是做到你能够解释：

> 为什么这么设计、如果服务挂了怎么办、如果消息重复怎么办、如果连接暴增怎么办、如果 Redis 挂了怎么办，以及为什么选择当前方案而不是其他方案。

这才是 40 天学习计划真正需要达到的目标。
