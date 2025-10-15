package main

import (
	"context"
	"fmt"
	"log"
	"sk_mall/utils"

	"github.com/elastic/go-elasticsearch/v7"
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9200",
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}
	res, err := client.Info()
	if err != nil {
		log.Fatalf("获取ES信息失败: %s", err)
	}
	defer res.Body.Close()
	r, e := utils.SearchDocuments(context.Background(), client, "sk_products", "te", 1, 10)
	if e != nil {
		fmt.Println("e:", e.Error())
		return
	}
	fmt.Println(*r)
}
