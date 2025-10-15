package utils

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"io"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type Document struct {
	EvID         string    `json:"ev_id"`
	ProductID    string    `json:"product_id"`
	Stock        int       `json:"stock"`
	StartTime    time.Time `json:"start_time"` // 直接使用 time.Time
	EndTime      time.Time `json:"end_time"`   // 直接使用 time.Time
	ProductName  string    `json:"product_name"`
	ProductDesc  string    `json:"product_desc"`
	ProductPrice float64   `json:"product_price"`
	CoverPath    string    `json:"coverpath"`
}

type SearchHit struct {
	Id     string   `json:"_id"`
	Source Document `json:"_source"`
}

type SearchResponse struct {
	Hits struct {
		Hits []SearchHit `json:"hits"`
	} `json:"hits"`
}

func SearchDocuments(ctx context.Context, client *elasticsearch.Client, indexName, query string, page, limit int) (ret *SearchResponse, err1 error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // 防止分页过大
	}

	from := (page - 1) * limit
	searchBody := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"product_name", "product_desc"},
			},
		},
		"from": from,
		"size": limit,
	}

	bodyBytes, err := json.Marshal(searchBody)
	if err != nil {
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{indexName},
		Body:  strings.NewReader(string(bodyBytes)),
	}

	resp, err := req.Do(ctx, client)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.IsError() {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(body))
	}

	// 解析响应
	var response SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}
