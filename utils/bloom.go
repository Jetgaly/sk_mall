package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisBloomFilter struct {
	client     *redis.Client
	key        string
	bitMapSize uint64                // 位图大小
	hashFuncs  []func([]byte) uint64 // 哈希函数
}
func (bf *RedisBloomFilter) AddToBloomFilterBatch(data string, batchSize int) error {
    if batchSize <= 0 {
        batchSize = 10 // 默认批次大小
    }
    
    // 计算所有位位置
    positions := make([]int64, 0, len(bf.hashFuncs))
    for _, hashFunc := range bf.hashFuncs {
        hashValue := hashFunc([]byte(data))
        bitPos := hashValue % bf.bitMapSize
        positions = append(positions, int64(bitPos))
    }
    
    // 分批处理
    for i := 0; i < len(positions); i += batchSize {
        end := i + batchSize
        if end > len(positions) {
            end = len(positions)
        }
        
        batch := positions[i:end]
        if err := bf.setBitsBatch(batch); err != nil {
            return err
        }
    }
    
    return nil
}

// 设置单个批次
func (bf *RedisBloomFilter) setBitsBatch(positions []int64) error {
    ctx := context.Background()
    pipe := bf.client.Pipeline()
    
    for _, pos := range positions {
        pipe.SetBit(ctx, bf.key, pos, 1)
    }
    
    _, err := pipe.Exec(ctx)
    return err
}

func (bf *RedisBloomFilter) ExistsInBloomFilterOptimized(data string) (bool, error) {
	ctx := context.Background()

	// 分批执行，遇到0就立即返回
	batchSize := 10 

	for i := 0; i < len(bf.hashFuncs); i += batchSize {
		end := i + batchSize
		if end > len(bf.hashFuncs) {
			end = len(bf.hashFuncs)
		}

		pipe := bf.client.Pipeline()
		cmds := make([]*redis.IntCmd, end-i)

		// 收集当前批次的命令
		for j := i; j < end; j++ {
			hashValue := bf.hashFuncs[j]([]byte(data))
			bitPos := hashValue % bf.bitMapSize
			cmds[j-i] = pipe.GetBit(ctx, bf.key, int64(bitPos))
		}

		// 执行当前批次
		if _, err := pipe.Exec(ctx); err != nil {
			return false, err
		}

		// 检查当前批次结果
		for _, cmd := range cmds {
			bit, err := cmd.Result()
			if err != nil {
				return false, err
			}
			if bit == 0 {
				return false, nil // 提前退出
			}
		}
	}

	return true, nil
}
