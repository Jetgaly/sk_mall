# skmall

## 简单架构图

![](https://raw.githubusercontent.com/Jetgaly/sk_mall/refs/heads/master/static/imgs/structure.png)

## desc

```
jjet@jet:~/projects/GoProjects/sk_mall$ tree -L 1
.
├── api 			#api服务
├── go.mod			
├── go.sum
├── readme.md
├── rpc				#rpc服务
├── scripts			
├── static
├── test
└── utils
```

该项目使用gozero微服务框架实现了秒杀商城的后端服务系统，主要实现了基本的功能接口和**核心的高并发商品秒杀系统**，旨在**高并发、高可用**，主要涉及以下内容：

- golang
- gozero
- dtm
- mysql
- redis
- canal
- rabbitmq
- elastic search
- docker
- etcd
- nginx
- zipkin
- prometheus
- grafana

1. gozero作为微服务的框架，
2. dtm协调**分布式事务**，保证**数据的一致性**，并且**异步调用**实现了服务的解耦，
3. redis作为**数据缓冲层**，解决**高并发**问题，还作为分布式锁的主要组件，实现**进程级别的分布式锁**，
4. mysql实现微服务的**分库管理**，避免共享大库，实现了**数据的解耦**，提高数据操作性能，以及增强**数据的安全性**，
5. canal通过监听mysql的binlog将商品数据同步到es，**减低微服务系统的耦合度**，
6. rabbitmq实现系统**核心异步化**，将订单超时取消、库存扣减消息等操作异步处理，提升**系统吞吐量**并实现模块间**解耦**。
7. es接入ik分词器，实现了商品的**模糊匹配**，
8. 使用docker进行各个服务的**快速部署**，
9. etcd实现**服务注册与发现**，
10. nginx对请求进行路由分发，实现**负载均衡**，
11. zipkin实现微服务请求链路的统计，以及可视化查询请求链路，在程序错误时可以根据日志的链路id进行**快速定位错误**，以及**数据恢复**，提升运维效率，
12. Prometheus，grafana实现了微服务指标的**统计和监控**，并且结合alertmanager实现相关的**预警报告**

## detail

### 商品、活动预热

在定时任务服务中实现了活动开始前的商品和活动的数据预热功能，将部分**数据预热缓存**到redis，有效地接下活动开始后的海量请求

### redis分布式锁

部署多个redis集群，并且设置保护协程来续费分布式锁，解决**分布式锁的误删**问题，使用redis分布式锁防止**缓存击穿**、**重复支付**等问题，

#### 随机expire

使用一定范围的随机数作为redis缓存数据的expire，有效防止**缓存雪崩**

### 核心秒杀功能

功能主要架构是：同步redis库存扣减，发送订单信息到rmq，再由order consumer来**异步消费**消息，操作mysql生成订单，**解耦**了秒杀接口直接操作mysql生成订单，从而实现该接口的高并发功能，性能瓶颈不再是mysql

- redis存储秒杀的库存信息，使用lua脚本和**乐观锁**解决库存扣减的原子性和**商品超卖问题**
- redis存储已经购买过的用户集合，解决**一人一单**问题
- 在order consumer订单落库时通过写入rmq实现的**延迟队列**来控制订单的**支付过期**状态，实现库存的归还

### 测试

在 **wsl2 8c16g** 的环境下**单机部署自压测试**(jmeter):`4500+ qps`

### 启动顺序：

1. product，
2. user，
3. merchant，
4. order_consumer，
5. seckill, order, 
6. aggr_order, 
7. payment ,
8. gateway

