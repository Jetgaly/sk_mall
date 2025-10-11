package main

import (
	
	"sk_mall/rpc/cron/cache_loader/tasks"
	
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	t := tasks.Scanner{}
	c := cron.New(cron.WithSeconds())
	// c.AddFunc("*/1 * * * * *", func() {
	// 	fmt.Println(strconv.Itoa(time.Now().Second()))
	// })
	c.AddJob("*/1 * * * * *", &t)
	c.Start()
	time.Sleep(10 * time.Second)
}
