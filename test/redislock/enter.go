package main

import (
	"context"
	"fmt"
	"sk_mall/utils"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
)

func t(args ...string) {
	for _, v := range args {
		fmt.Println(v)
	}
}

func main() {
	r1 := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "Qwe@123456",
	})
	r2 := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6389",
		Password: "123456",
	})
	r3 := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6399",
		Password: "123456",
	})
	t, _ := utils.NewRedLockCreater([]*redis.Client{r1, r2, r3})
	ctx, cancel := context.WithCancel(context.Background())
	m, e := t.GetLock(ctx, "testtest",redsync.WithTries(32))
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	time.Sleep(35 * time.Second)
	t.ReleaseLock(m,cancel)
}

// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/go-redsync/redsync/v4"
// 	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
// 	"github.com/redis/go-redis/v9"
// )

// func main() {
// 	// 多个 redis 节点
// 	pools := []redsync.Pool{
// 		goredis.NewPool(redis.NewClient(&redis.Options{Addr: "localhost:6379"})),
// 		goredis.NewPool(redis.NewClient(&redis.Options{Addr: "localhost:6380"})),
// 		goredis.NewPool(redis.NewClient(&redis.Options{Addr: "localhost:6381"})),
// 	}

// 	// 创建 redsync 实例
// 	rs := redsync.New(pools...)

// 	// 创建互斥锁
// 	mutex := rs.NewMutex("my-global-lock",
// 		redsync.WithExpiry(10*time.Second), // 锁过期时间
// 		redsync.WithTries(3),               // 获取锁重试次数
// 	)

// 	// 获取锁
// 	if err := mutex.Lock(); err != nil {
// 		fmt.Println("加锁失败:", err)
// 		return
// 	}
// 	fmt.Println("加锁成功")

// 	// 模拟业务
// 	time.Sleep(5 * time.Second)

// 	// 释放锁
// 	if ok, err := mutex.Unlock(); !ok || err != nil {
// 		fmt.Println("释放锁失败:", err)
// 	} else {
// 		fmt.Println("释放锁成功")
// 	}
// }
